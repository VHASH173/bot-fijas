package comandos

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GenerarPlanes arma el menú de precios y suscripciones VIP
func GenerarPlanes() (string, tgbotapi.InlineKeyboardMarkup) {
	caption := `<b>❰ #DORADO FIJAS ❱ ➢ PLANES DISPONIBLES</b>

<b>[ PACKS DE PRONÓSTICOS ]</b>

🔰 <b>BÁSICO</b>
➔ 3 Fijas VIP | S/ 10
➔ 5 Fijas VIP + 1 Combinada | S/ 15
➔ 10 Fijas VIP + 3 Combinadas | S/ 25

⭐ <b>STANDARD</b>
➔ 15 Fijas VIP + 5 Combinadas | S/ 40
➔ 20 Fijas VIP + Bankroll Mngt | S/ 55

💎 <b>PREMIUM</b>
➔ 30 Fijas VIP + Asesoría 1a1 | S/ 80
➔ 50 Fijas VIP + Reto Escalera | S/ 120

━━━━━━━━━━━━━━━━━━━━━━━━

<b>[ SUSCRIPCIONES VIP POR DÍAS ]</b>

🔰 <b>BÁSICO</b>
➔ 3 días VIP | S/ 15
➔ 7 días VIP | S/ 25
➔ 10 días VIP | S/ 35

⭐ <b>STANDARD</b>
➔ 15 días VIP | S/ 45
➔ 30 días VIP | S/ 70
➔ 45 días VIP | S/ 95

💎 <b>PREMIUM</b>
➔ 60 días VIP | S/ 120
➔ 90 días VIP | S/ 160
➔ 120 días VIP | S/ 200`

	// Botonera 2x2 calcada de tu referencia
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("【💬】 GRUPO VIP", "https://t.me/TuGrupoAqui"),
			tgbotapi.NewInlineKeyboardButtonURL("【♻️】 SOPORTE", "https://t.me/@LaRealFijaVIP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("【👑】 FUNDADOR", "https://t.me/@LaRealFijaVIP"),
			tgbotapi.NewInlineKeyboardButtonURL("【🛒】 COMPRAR", "https://t.me/@LaRealFijaVIP"), // Link directo a ti para cerrar la venta
		),
	)

	return caption, keyboard
}
