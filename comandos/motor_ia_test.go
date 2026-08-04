package comandos

import "testing"

func TestConstruirDatosTicketUsaPartidoRealYPronostico(t *testing.T) {
	ticket := construirDatosTicket("Liga Premier", "Barcelona", "Real Madrid", "1X2", "Barcelona a ganar")

	if ticket.Partido != "Barcelona vs Real Madrid" {
		t.Fatalf("se esperaba el partido real, obtuve %q", ticket.Partido)
	}
	if ticket.Mercado != "1X2" {
		t.Fatalf("se esperaba el mercado 1X2, obtuve %q", ticket.Mercado)
	}
	if ticket.Pronostico != "Barcelona a ganar" {
		t.Fatalf("se esperaba el pronóstico real, obtuve %q", ticket.Pronostico)
	}
	if ticket.Cuota != "Sin datos" {
		t.Fatalf("se esperaba el valor por defecto de cuota, obtuve %q", ticket.Cuota)
	}
}
