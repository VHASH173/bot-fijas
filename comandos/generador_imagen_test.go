package comandos

import (
	"strings"
	"testing"
)

func TestGenerarHTMLTicket(t *testing.T) {
	datos := DatosTicket{
		Partido:    "Real Madrid vs Man City",
		Mercado:    "Ganador del partido",
		Pronostico: "Real Madrid a ganar",
		Cuota:      "2.85",
	}

	html := GenerarHTMLTicket(datos)
	for _, want := range []string{"Real Madrid vs Man City", "Ganador del partido", "Real Madrid a ganar", "2.85"} {
		if !strings.Contains(html, want) {
			t.Fatalf("el HTML generado no contiene %q\n%s", want, html)
		}
	}
}
