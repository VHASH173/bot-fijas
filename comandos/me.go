package comandos

import (
	"fmt"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GenerarMe(user *tgbotapi.User, fechaRegistro time.Time, consultas int, consultasHoy int) string {
	userName := user.UserName
	if userName == "" {
		userName = "No_disponible"
	}

	diasActivo := int(time.Since(fechaRegistro).Hours() / 24)
	if diasActivo == 0 {
		diasActivo = 1 
	}

	caption := fmt.Sprintf(`<b>[BOT DE LAS JIJAS] ➾ ME - PERFIL</b>

<b>PERFIL DE ➾</b> %s

<b>INFORMACIÓN PERSONAL</b>

[🙎‍♂️] ID ➾ <code>%d</code>
[👨🏻‍💻] USER ➾ @%s
[👺] ESTADO ➾ ACTIVO
[📅] F. REGISTRO ➾ %s

<b>ESTADO DE CUENTA</b>

[〽️] ROL ➾ USUARIO
[⏱️] ANTI-SPAM ➾ 0'
[⏳] TIEMPO ➾ %d DÍAS
[📅] F. EXPIRACION ➾ 3000-01-01

<b>USO DEL SERVICIO</b>

[📊] CONSULTAS ➾ %d
[📅] CONSULTAS DE HOY ➾ %d`,
		user.FirstName,
		user.ID,
		userName,
		fechaRegistro.Format("2006-01-02"),
		diasActivo,
		consultas,
		consultasHoy,
	)

	return caption
}