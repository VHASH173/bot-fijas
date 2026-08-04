package comandos

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GenerarMenuCmds(categoria string) (string, tgbotapi.InlineKeyboardMarkup) {
	var caption string
	var keyboard tgbotapi.InlineKeyboardMarkup

	botonVolver := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("← Volver al menú", "btn_volver"),
	)

	switch categoria {
	case "btn_fija":
		// AHORA ESTE ES EL SUBMENÚ DE MERCADOS
		caption = `<b>[BOT DE LAS JIJAS] ➢ SELECCIONA EL MERCADO</b>

Nuestra IA analiza más de 40 ligas en tiempo real. 
Selecciona qué tipo de fija estás buscando hoy:`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【🏆】 1X2 (GANADOR)", "ia_1x2"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【🥅】 GOLES (O/U)", "ia_goles"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【🔥】 AMBOS ANOTAN", "ia_btts"),
			),
			botonVolver,
		)

	case "btn_vip":
		caption = `<b>[BOT DE LAS JIJAS]</b> ➞ SISTEMA PREDICTIVO

<b>CATEGORÍA</b> ➞ VIP Premium
<b>ESTADO</b> ➞ SOLO MIEMBROS [💎]
<b>MERCADO</b> ➞ Pre-Match Combinadas

[💎] MOTOR MATEMÁTICO VIP

<b>CUOTA APROX</b> ➞ 3.00 - 5.00+
<b>WINRATE VIP</b> ➞ 82% 
<b>ACCESO</b> ➞ DENEGADO 🚫

<i>Adquiere tu membresía para desbloquear los pronósticos combinados de máxima rentabilidad.</i>`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "btn_stats":
		caption = `<b>[BOT DE LAS JIJAS]</b> ➞ SISTEMA PREDICTIVO

<b>CATEGORÍA</b> ➞ Estadísticas del Mes
<b>ESTADO</b> ➞ OPERATIVO [✅]

[📊] HISTORIAL DE ACIERTOS

<b>FIJAS ACERTADAS</b> ➞ 45
<b>FIJAS FALLADAS</b> ➞ 12
<b>YIELD MENSUAL</b> ➞ +18.5%

<i>Transparencia total. El modo Jijas no perdona a las casas de apuestas.</i>`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "btn_planes":
		caption = `<b>[BOT DE LAS JIJAS]</b> ➞ SUSCRIPCIONES VIP

<b>CATEGORÍA</b> ➞ Planes y Precios
<b>ESTADO</b> ➞ DISPONIBLE [✅]

[💳] PLANES DISPONIBLES:

🥇 <b>PLAN SEMANAL:</b> Acceso al motor combinadas (7 días)
🏆 <b>PLAN MENSUAL:</b> Acceso total + Gestión de Bank (30 días)

<i>Contacta al administrador para los métodos de pago.</i>`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "cmd_referidos":
		caption = `<b>【 #KING DATA 】 ➤ SISTEMA DE REFERIDOS</b>

<b>USUARIO</b> ➤ El Dorado 👑

<b>ESTADÍSTICAS GENERALES</b>
• Total invitados: 0
• Ganancias: 0 créditos

<b>LISTA DE REGISTROS</b>
No hay registros disponibles.

<b>ENLACE DE INVITACIÓN</b>
https://t.me/KingDataX_bot?start=7010388601`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "cmd_compras":
		caption = `<b>【 #KING DATA 】 ➤ HISTORIAL DE COMPRAS</b>

No tienes registros de compras o activaciones.`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "cmd_historial_fijas":
		caption = `<b>【 #KING DATA 】 ➤ HISTORIAL VACIÓ</b>

⚠️ No hay consultas registradas para este usuario.`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	case "cmd_suscripcion":
		caption = `<b>【 #KING DATA 】 ➤ SUSCRIPCIÓN</b>

Tu suscripción actual es gratuita.

Actualiza a VIP para obtener acceso completo a todas las fijas y beneficios.`
		keyboard = tgbotapi.NewInlineKeyboardMarkup(botonVolver)

	default: // Menú Principal
		caption = "<b>[BOT DE LAS JIJAS]</b> → <i>Menú Principal</i>\n\nSelecciona una categoría de apuestas:"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【⚽】 SOLTAR FIJA", "btn_fija"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【💎】 ZONA VIP", "btn_vip"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【📊】 ESTADÍSTICAS", "btn_stats"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("【💳】 PLANES VIP", "btn_planes"),
			),
		)
	}

	return caption, keyboard
}
