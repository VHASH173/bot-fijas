package comandos

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GenerarMe(user *tgbotapi.User, fechaRegistro time.Time, rol, plan string, creditos int, esVIP bool, fechaVencimientoVIP time.Time, consultasTotales int, consultasHoy int, botUsername string) (string, tgbotapi.InlineKeyboardMarkup) {
	fullName := strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	if fullName == "" {
		fullName = "No disponible"
	}

	username := "No disponible"
	if user.UserName != "" {
		username = fmt.Sprintf("@%s", user.UserName)
	}

	fechaRegistroTexto := fechaRegistro.Format("02/01/2006 15:04:05")
	if fechaRegistro.IsZero() {
		fechaRegistroTexto = "N/D"
	}

	diasRestantes := "N/A"
	fechaVencimientoTexto := "N/A"
	if esVIP {
		dias := int(time.Until(fechaVencimientoVIP).Hours() / 24)
		if dias < 0 {
			dias = 0
		}
		diasRestantes = fmt.Sprintf("%d DÍAS", dias)
		fechaVencimientoTexto = fechaVencimientoVIP.Format("02/01/2006")
	}

	caption := fmt.Sprintf(`<b>❰ #BOT_FIJAS ❱ ➣ PERFIL DE USUARIO</b>

「🆔」 • ID ➣ <code>%d</code>
「🙎」 • NOMBRE ➣ %s
「👨🏻‍💻」 • USUARIO ➣ %s
「✅️」 • ESTADO ➣ ACTIVO
「📅」 • F.REGISTRO ➣ %s

<b>💳 SUSCRIPCIÓN</b>

「〽️」 • ROL ➣ %s
「📈」 • PLAN ➣ %s
「⏱️」 • DIAS ➣ %s
「💰」 • CREDITOS ➣ %d
「⏳」 • VENCE ➣ %s

<b>📈 ACTIVIDAD Y USO</b>

「📊」 • CONSULTAS ➣ %d
「📅」 • CONSULTAS HOY ➣ %d
「⏱️」 • ANTI-SPAM ➣ 0 seg.

🔗 <b>ENLACE DE REFERIDO</b>
https://t.me/%s?start=%d`, user.ID, fullName, username, fechaRegistroTexto, rol, plan, diasRestantes, creditos, fechaVencimientoTexto, consultasTotales, consultasHoy, botUsername, user.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("👥 REFERIDOS", "cmd_referidos")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🛒 HISTORIAL COMPRAS", "cmd_compras")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📜 HISTORIAL APUESTAS", "cmd_historial")),
	)

	return caption, keyboard
}
