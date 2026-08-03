package comandos

import (
	"strings"
	"testing"
)

func TestGenerarHTMLTicket(t *testing.T) {
	datos := DatosTicket{
		Liga:       "Champions League",
		Local:      "Real Madrid",
		Visita:     "Man City",
		Mercado:    "Ganador del partido",
		Pronostico: "Real Madrid a ganar",
		Cuota:      "2.85",
		Stake:      "8",
	}

	html, err := GenerarHTMLTicket(datos)
	if err != nil {
		t.Fatalf("GenerarHTMLTicket devolvió error: %v", err)
	}

	for _, want := range []string{"Champions League", "Real Madrid", "Man City", "2.85", "8"} {
		if !strings.Contains(html, want) {
			t.Fatalf("el HTML generado no contiene %q\n%s", want, html)
		}
	}
}
