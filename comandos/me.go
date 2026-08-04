package comandos

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GenerarMe(user *tgbotapi.User, fechaRegistro time.Time, rol, plan string, creditos int, esVIP bool, fechaVencimientoVIP time.Time, consultasTotales int, consultasHoy int, botUsername string) string {
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

	rolUsuario := rol
	if rolUsuario == "" {
		rolUsuario = "Free"
	}

	planActual := plan
	if planActual == "" {
		planActual = "GRATIS"
	}

	diasRestantes := "0"
	fechaVencimientoTexto := "N/A"
	if esVIP {
		dias := int(time.Until(fechaVencimientoVIP).Hours() / 24)
		if dias < 0 {
			dias = 0
		}
		diasRestantes = fmt.Sprintf("%d", dias)
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
https://t.me/%s?start=%d

📌 <b>MENÚ DE OPCIONES</b>
「👥」 • REFERIDOS ➣ /referidos
「🛒」 • COMPRAS ➣ /compras
「📜」 • HISTORIAL ➣ /historial`, user.ID, fullName, username, fechaRegistroTexto, rolUsuario, planActual, diasRestantes, creditos, fechaVencimientoTexto, consultasTotales, consultasHoy, botUsername, user.ID)

	return caption
}
