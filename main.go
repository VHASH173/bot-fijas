package main

import (
	"log"
	"os" // Importamos OS para leer las variables de Railway
	"time"

	"bot-fijas/comandos" // Asegúrate de que coincida con el nombre de tu módulo

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ImagePath = "image.png" // Tu foto (asegúrate de que sea cuadrada 1:1)

type Perfil struct {
	FechaRegistro    time.Time
	ConsultasTotales int
	ConsultasHoy     int
	UltimaConsulta   time.Time // Controla el spam de 2 minutos
	EsVIP            bool      // Controla si tiene límite de 3 o es ilimitado
}

var baseDeDatos = make(map[int64]*Perfil)

func main() {
	// Leemos el token desde las variables de Railway
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Falta la variable TELEGRAM_BOT_TOKEN en el servidor")
	}

	// Iniciamos el bot con ese token
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Error al iniciar el bot: ", err)
	}
	bot.Debug = false
	log.Printf("¡Bot encendido! Logueado como %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var userID int64
		var chatID int64
		var tUser *tgbotapi.User

		if update.Message != nil {
			tUser = update.Message.From
			userID = tUser.ID
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			tUser = update.CallbackQuery.From
			userID = tUser.ID
			chatID = update.CallbackQuery.Message.Chat.ID
		} else {
			continue
		}

		// Registro de usuario
		if _, existe := baseDeDatos[userID]; !existe {
			baseDeDatos[userID] = &Perfil{
				FechaRegistro:    time.Now(),
				ConsultasTotales: 0,
				ConsultasHoy:     0,
				UltimaConsulta:   time.Now().Add(-10 * time.Minute), // Para que pueda usarlo apenas entra
				EsVIP:            false,                             // Por defecto todos son Free
			}
		}

		// Declaramos la variable usuario para usarla abajo
		usuario := baseDeDatos[userID]

		// --- MANEJO DE COMANDOS DE TEXTO ---
		if update.Message != nil && update.Message.IsCommand() {
			msg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(ImagePath))
			msg.ParseMode = "HTML"

			switch update.Message.Command() {
			case "start":
				usuario.ConsultasTotales++
				caption, keyboard := comandos.GenerarStart(tUser)
				msg.Caption = caption
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "register":
				usuario.ConsultasTotales++
				texto := "Registro completado"
				textMsg := tgbotapi.NewMessage(chatID, texto)
				textMsg.ParseMode = "HTML"
				bot.Send(textMsg)

			case "me":
				usuario.ConsultasTotales++
				msg.Caption = "Perfil de usuario"
				bot.Send(msg)

			case "cmds":
				usuario.ConsultasTotales++
				caption, keyboard := comandos.GenerarMenuCmds("menu_principal")
				msg.Caption = caption
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "planes":
				usuario.ConsultasTotales++
				caption, keyboard := comandos.GenerarMenuCmds("btn_planes")
				msg.Caption = caption
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "staff":
				usuario.ConsultasTotales++
				texto := "Contacto Staff"
				textMsg := tgbotapi.NewMessage(chatID, texto)
				textMsg.ParseMode = "HTML"
				bot.Send(textMsg)
			}
			continue
		}

		// --- MANEJO DE BOTONES (Callback) ---
		if update.CallbackQuery != nil {
			bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			messageID := update.CallbackQuery.Message.MessageID

			var caption string
			var keyboard tgbotapi.InlineKeyboardMarkup
			data := update.CallbackQuery.Data

			// AQUI CONECTAMOS LOS NUEVOS BOTONES DE LA IA
			if data == "ia_1x2" || data == "ia_goles" || data == "ia_btts" {
				var seCobroIntento bool
				caption, keyboard, seCobroIntento, usuario.ConsultasHoy, usuario.UltimaConsulta = comandos.GenerarPronosticoIA(usuario.ConsultasHoy, usuario.UltimaConsulta, usuario.EsVIP, data)

				if seCobroIntento {
					usuario.ConsultasTotales++ // Aumentamos las totales solo si la IA le dio respuesta

					// Ejemplo de uso para enviar un ticket como imagen en Telegram
					datosTicket := comandos.DatosTicket{
						Liga:       "Champions League",
						Local:      "Real Madrid",
						Visita:     "Man City",
						Mercado:    "Ganador del partido",
						Pronostico: "Real Madrid a ganar",
						Cuota:      "2.85",
						Stake:      "8",
					}
					htmlTicket, err := comandos.GenerarHTMLTicket(datosTicket)
					if err == nil {
						urlImagen, errImg := comandos.ConvertirHTMLaImagen(htmlTicket)
						if errImg == nil && urlImagen != "" {
							mensajeFoto := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(urlImagen))
							mensajeFoto.Caption = "🔥 ¡NUEVA FIJA CONFIRMADA! 🔥\n\nRegístrate con nuestro código en Betano para seguir este pronóstico."
							bot.Send(mensajeFoto)
						}
					}
				}
			} else if data == "btn_volver" {
				caption, keyboard = comandos.GenerarMenuCmds("menu_principal")
			} else {
				caption, keyboard = comandos.GenerarMenuCmds(data)
			}

			editMsg := tgbotapi.NewEditMessageCaption(chatID, messageID, caption)
			editMsg.ParseMode = "HTML"
			editMsg.ReplyMarkup = &keyboard
			bot.Send(editMsg)
		}
	}
}
