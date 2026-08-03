package comandos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func GenerarPronosticoIA(consultasHoy int, ultimaConsulta time.Time, esVIP bool, mercado string) (string, tgbotapi.InlineKeyboardMarkup, bool, int, time.Time) {

	oddsAPIKey := os.Getenv("ODDS_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if oddsAPIKey == "" || geminiKey == "" {
		log.Println("Faltan variables de entorno: ODDS_API_KEY o GEMINI_API_KEY")
		teclado := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER 】", "btn_volver")),
		)
		return "⚠️ No tengo configuradas las claves de IA y apuestas. Activa las variables de entorno.", teclado, false, consultasHoy, ultimaConsulta
	}

	// 1. FILTRO ANTI-SPAM (2 Minutos)
	tiempoPasado := time.Since(ultimaConsulta)
	if tiempoPasado < 2*time.Minute {
		faltan := 2*time.Minute - tiempoPasado
		minutos := int(faltan.Minutes())
		segundos := int(faltan.Seconds()) % 60

		textoSpam := fmt.Sprintf("⏳ <b>¡Epa, tranquilo máquina!</b>\n\nEl algoritmo está procesando datos. Debes esperar <b>%d min y %d seg</b> para solicitar otra fija.", minutos, segundos)
		teclado := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER 】", "btn_volver")))
		return textoSpam, teclado, false, consultasHoy, ultimaConsulta
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
		return textoLimite, tecladoVIP, false, consultasHoy, ultimaConsulta
	}

	// 3. OBTENER PARTIDO REAL Y CUOTAS (The Odds API)
	// Buscamos los próximos partidos de fútbol (soccer) a nivel mundial
	url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/upcoming/odds/?apiKey=%s&regions=eu&markets=h2h", oddsAPIKey)
	respOdds, errOdds := http.Get(url)

	var liga, local, visita, infoCuotas string

	if errOdds != nil || respOdds.StatusCode != 200 {
		log.Println("Error conectando a The Odds API")
		liga, local, visita, infoCuotas = "Fútbol Global", "Equipo A", "Equipo B", "Cuotas no disponibles momentáneamente"
	} else {
		defer respOdds.Body.Close()
		var eventos []EventoDeportivo
		json.NewDecoder(respOdds.Body).Decode(&eventos)

		if len(eventos) > 0 {
			evento := eventos[0] // Tomamos el partido más próximo
			liga = evento.SportTitle
			local = evento.HomeTeam
			visita = evento.AwayTeam

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
		prompt := fmt.Sprintf(`Eres un tipster profesional de apuestas deportivas muy agresivo y persuasivo. 
Tengo este partido REAL a punto de jugarse:
Liga: %s
Partido: %s vs %s
Cuotas Actuales: %s

El usuario está buscando un pronóstico para el mercado: %s.
Basado en las cuotas, genera un pronóstico de máximo 4 líneas. Sé muy seguro de ti mismo, usa un lenguaje de "fija" o "verde seguro", añade emojis de dinero o fuego y recomienda un stake del 1 al 10. No pongas asteriscos ni formato markdown extra, solo el texto limpio.`, liga, local, visita, infoCuotas, mercado)

		resp, err := model.GenerateContent(ctx, genai.Text(prompt))

		if err != nil {
			log.Println("Error generando contenido:", err)
			respuestaIA = "⚠️ La IA está sobrecargada analizando métricas de este partido. Intenta en unos minutos."
		} else {
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						respuestaIA = fmt.Sprintf("%v", part)
					}
				}
			}
		}
	}

	// 5. ARMAMOS EL MENSAJE FINAL PARA TELEGRAM
	textoExito := fmt.Sprintf(`<b>[BOT DE LAS JIJAS] ➢ ANÁLISIS EN VIVO 🤖</b>

🏆 <b>Competición:</b> %s
⚽ <b>Encuentro:</b> %s vs %s
🎯 <b>Mercado analizado:</b> %s

<b>[🔥] PRONÓSTICO DEL SISTEMA:</b>
%s

<i>⚠️ Invierte con cabeza. Las cuotas varían según tu casa de apuestas.</i>`, liga, local, visita, mercado, respuestaIA)

	teclado := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("【 🏠 VOLVER AL MENÚ 】", "btn_volver")))

	// Todo salió bien, cobramos el intento y actualizamos la hora
	return textoExito, teclado, true, consultasHoy + 1, time.Now()
}
