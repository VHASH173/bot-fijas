package comandos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Estructuras para leer los datos reales de The-Odds-API
type EventoDeportivo struct {
	SportTitle string   `json:"sport_title"`
	HomeTeam   string   `json:"home_team"`
	AwayTeam   string   `json:"away_team"`
	Bookmakers []Bookie `json:"bookmakers"`
}

type Bookie struct {
	Title   string   `json:"title"`
	Markets []Market `json:"markets"`
}

type Market struct {
	Key      string    `json:"key"`
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func construirDatosTicket(liga, local, visita, mercado, pronostico string) DatosTicket {
	partido := "Partido en análisis"
	local = strings.TrimSpace(local)
	visita = strings.TrimSpace(visita)
	if local != "" || visita != "" {
		partido = fmt.Sprintf("%s vs %s", local, visita)
	}

	mercadoTexto := strings.TrimSpace(mercado)
	if mercadoTexto == "" {
		mercadoTexto = "Mercado principal"
	}

	pronosticoTexto := strings.TrimSpace(pronostico)
	if pronosticoTexto == "" {
		pronosticoTexto = "Pronóstico en análisis"
	}

	return DatosTicket{
		Partido:    partido,
		Mercado:    mercadoTexto,
		Pronostico: pronosticoTexto,
		Cuota:      "Sin datos",
	}
}

func respuestaIAValida(respuesta string) bool {
	texto := strings.TrimSpace(strings.ToLower(respuesta))
	if texto == "" {
		return false
	}
	if strings.Contains(texto, "error") || strings.Contains(texto, "sobrecargada") || strings.Contains(texto, "intenta") || strings.Contains(texto, "no tengo") {
		return false
	}
	return true
}

func GenerarPronosticoIA(consultasHoy int, ultimaConsulta time.Time, esVIP bool, mercado string) (string, tgbotapi.InlineKeyboardMarkup, DatosTicket, bool, bool, int, time.Time) {

	oddsAPIKey := os.Getenv("ODDS_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if oddsAPIKey == "" || geminiKey == "" {
		log.Println("Faltan variables de entorno: ODDS_API_KEY o GEMINI_API_KEY")
		teclado := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER 】", "btn_volver")),
		)
		return "⚠️ No tengo configuradas las claves de IA y apuestas. Activa las variables de entorno.", teclado, DatosTicket{}, false, false, consultasHoy, ultimaConsulta
	}

	// 1. FILTRO ANTI-SPAM (2 Minutos)
	tiempoPasado := time.Since(ultimaConsulta)
	if tiempoPasado < 2*time.Minute {
		faltan := 2*time.Minute - tiempoPasado
		minutos := int(faltan.Minutes())
		segundos := int(faltan.Seconds()) % 60

		textoSpam := fmt.Sprintf("⏳ <b>¡Epa, tranquilo máquina!</b>\n\nEl algoritmo está procesando datos. Debes esperar <b>%d min y %d seg</b> para solicitar otra fija.", minutos, segundos)
		teclado := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER 】", "btn_volver")))
		return textoSpam, teclado, DatosTicket{}, false, false, consultasHoy, ultimaConsulta
	}

	// 2. LÍMITE FREEMIUM (Max 3 al día)
	if !esVIP && consultasHoy >= 3 {
		textoLimite := `🚫 <b>¡LÍMITE DIARIO ALCANZADO!</b> 🚫

Has consumido tus 3 fijas gratuitas de hoy. La IA ha detectado oportunidades de alta rentabilidad para las próximas horas.

💎 <b>Adquiere un Plan VIP para soltar pronósticos ilimitados y reventar el bank.</b>`

		tecladoVIP := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【💳】 VER PLANES VIP", "btn_planes")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER 】", "btn_volver")),
		)
		return textoLimite, tecladoVIP, DatosTicket{}, false, false, consultasHoy, ultimaConsulta
	}

	// 3. OBTENER PARTIDO REAL Y CUOTAS (The Odds API)
	// Buscamos los próximos partidos de fútbol (soccer) a nivel mundial
	url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/upcoming/odds/?apiKey=%s&regions=eu&markets=h2h", oddsAPIKey)
	respOdds, errOdds := http.Get(url)

	var liga, local, visita, infoCuotas, bookmaker string

	if errOdds != nil || respOdds.StatusCode != 200 {
		log.Println("Error conectando a The Odds API")
		liga, local, visita, infoCuotas, bookmaker = "Fútbol Global", "Equipo A", "Equipo B", "Cuotas no disponibles momentáneamente", "Sin datos"
	} else {
		defer respOdds.Body.Close()
		var eventos []EventoDeportivo
		json.NewDecoder(respOdds.Body).Decode(&eventos)

		if len(eventos) > 0 {
			evento := eventos[0] // Tomamos el partido más próximo
			liga = evento.SportTitle
			local = evento.HomeTeam
			visita = evento.AwayTeam

			if len(evento.Bookmakers) > 0 {
				bookmaker = evento.Bookmakers[0].Title
			}

			// Extraemos las cuotas si el bookmaker las liberó
			if len(evento.Bookmakers) > 0 && len(evento.Bookmakers[0].Markets) > 0 {
				outcomes := evento.Bookmakers[0].Markets[0].Outcomes
				infoCuotas = ""
				for _, out := range outcomes {
					infoCuotas += fmt.Sprintf("%s: %.2f | ", out.Name, out.Price)
				}
			} else {
				infoCuotas = "Cuotas ocultas por la casa"
			}
		}
	}

	// 4. CONEXIÓN A GEMINI IA (Enviando los datos reales)
	ctx := context.Background()
	client, errGenai := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	var respuestaIA string

	if errGenai != nil {
		log.Println("Error conectando a Gemini:", errGenai)
		respuestaIA = "❌ Error en los servidores de IA. Intenta de nuevo."
	} else {
		defer client.Close()
		model := client.GenerativeModel("gemini-1.5-flash")

		// El "Prompt" maestro que mezcla la IA con los datos reales
		prompt := fmt.Sprintf(`Eres un analista de apuestas. Responde solo con este formato exacto en Markdown, sin saludos, sin introducciones, sin despedidas:

🔥 **APUESTA CONFIRMADA** 🔥
⚽ **Partido:** [Equipo A vs Equipo B]
🎯 **Mercado:** [Ej: Ganador del partido / Más de 2.5 goles]
📊 **Pronóstico:** [El pick exacto]
📈 **Cuota:** [Ej: 1.85]
🏦 **Casa de Apuestas:** [Ej: Betano, Bet365, 1xBet]
💰 **Stake Recomendado:** [Ej: Stake 1.5 o Monto sugerido]

Usa estos datos:
Liga: %s
Partido: %s vs %s
Bookmaker: %s
Cuotas Actuales: %s
Mercado: %s`, liga, local, visita, bookmaker, infoCuotas, mercado)

		resp, err := model.GenerateContent(ctx, genai.Text(prompt))

		if err != nil {
			log.Println("Error generando contenido:", err)
			respuestaIA = "⚠️ La IA está analizando demasiadas métricas. Intenta en unos minutos."
		} else {
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						respuestaIA = fmt.Sprintf("%v", part)
					}
				}
			}
			if !respuestaIAValida(respuestaIA) {
				respuestaIA = "⚠️ La IA está analizando demasiadas métricas. Intenta en unos minutos."
			}
		}
	}

	// 5. ARMAMOS EL MENSAJE FINAL PARA TELEGRAM
	textoExito := fmt.Sprintf(`<b>[BOT DE LAS JIJAS] ➢ ANÁLISIS EN VIVO 🤖</b>

🏆 <b>Competición:</b> %s
⚽ <b>Encuentro:</b> %s vs %s
🎯 <b>Mercado analizado:</b> %s
🏦 <b>Bookmaker:</b> %s
📊 <b>Cuotas:</b> %s

<b>[🔥] PRONÓSTICO DEL SISTEMA:</b>
%s

<i>⚠️ Invierte con cabeza. Las cuotas varían según tu casa de apuestas.</i>`, liga, local, visita, mercado, bookmaker, infoCuotas, respuestaIA)

	teclado := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER AL MENÚ 】", "btn_volver")))

	debeEnviarImagen := respuestaIAValida(respuestaIA)
	if !debeEnviarImagen {
		log.Println("Se omite la imagen porque la respuesta de IA no fue válida")
	}

	ticket := construirDatosTicket(liga, local, visita, mercado, respuestaIA)

	// Todo salió bien, cobramos el intento y actualizamos la hora
	return textoExito, teclado, ticket, true, debeEnviarImagen, consultasHoy + 1, time.Now()
}
