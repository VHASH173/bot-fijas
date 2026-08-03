package comandos

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GenerarStart arma el mensaje de bienvenida con la arquitectura VIP de bloques
func GenerarStart(user *tgbotapi.User) (string, tgbotapi.InlineKeyboardMarkup) {
	// Mención en HTML (crea el link al perfil del usuario)
	mencion := fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", user.ID, user.FirstName)

	// Generamos la fecha actual automáticamente
	fechaActual := time.Now().Format("02-01-2006")

	// Usamos líneas continuas para recrear el separador del diseño
	linea := "━━━━━━━━━━━━━━━━━━━━━━━━"

	caption := fmt.Sprintf(`Hola, %s 

Estás usando: <b>BOT DE LAS JIJAS</b>
Un motor diseñado para reventar a las casas de apuestas con precisión y rentabilidad.
%s
<b>[ 📍 ] DETALLES DEL SISTEMA</b>
<b>Usuario</b> ➢ @DoradoFijas_bot
<b>Versión</b> ➢ v2.0 VIP
<b>Fecha</b> ➢ %s
%s
<b>[ 📍 ] COMANDOS PRINCIPALES</b>
/register ➢ Registro de cuenta VIP
/cmds ➢ Menú de apuestas y fijas
/me ➢ Perfil y bankroll del usuario
/planes ➢ Suscripciones VIP
/staff ➢ Contacto y soporte directo
%s
<b>[ 📍 ] AVISO</b>
<i>Gestión de bank obligatoria. Las apuestas conllevan riesgo, opera con responsabilidad.</i>
%s
<i>[ 📍 ] Servicio administrado por EL FUNDADOR</i>`, mencion, linea, fechaActual, linea, linea, linea)

	// Botones en cuadrícula 2x2 (igual que en tu captura de referencia)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("【💬】 CANAL VIP", "https://t.me/TuGrupoAqui"),
			tgbotapi.NewInlineKeyboardButtonURL("【♻️】 SOPORTE", "https://t.me/TuSoporteAqui"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("【👑】 FUNDADOR", "https://t.me/TuUsuarioAqui"),
			tgbotapi.NewInlineKeyboardButtonURL("【📊】 RESULTADOS", "https://t.me/TuCanalResultados"),
		),
	)

	return caption, keyboard
}
