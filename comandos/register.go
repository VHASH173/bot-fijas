package comandos

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GenerarRegister arma la respuesta visual de registro
func GenerarRegister(user *tgbotapi.User) string {
	// Verificamos si el usuario tiene un @username configurado en Telegram
	username := "Sin @Usuario"
	if user.UserName != "" {
		username = "@" + user.UserName
	}

	// Mención en HTML
	mencion := fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", user.ID, user.FirstName)

	// Estructura exacta a la de tu referencia visual
	texto := fmt.Sprintf(`<b>❰ #DORADO FIJAS ❱ ➢ REGISTRO COMPLETADO</b>

♻️ Tus datos fueron actualizados correctamente.

👤 <b>USUARIO:</b> %s
📌 <b>USER:</b> %s
🆔 <b>ID:</b> <code>%d</code>

📋 <i>Usa /me para ver tu perfil actualizado.</i>`, mencion, username, user.ID)

	return texto
}
