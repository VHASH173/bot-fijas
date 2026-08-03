package comandos

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DatosTicket almacena la info que inyectaremos en el HTML
type DatosTicket struct {
	Liga       string
	Local      string
	Visita     string
	Mercado    string
	Pronostico string
	Cuota      string
	Stake      string
}

// Estructura para leer la respuesta de la API HCTI
type HCTIResponse struct {
	URL string `json:"url"`
}

// GenerarHTMLTicket toma los datos y los mete en la plantilla HTML
func GenerarHTMLTicket(datos DatosTicket) (string, error) {
	templatePath := filepath.Join("plantillas", "ticket.html")
	if _, err := os.Stat(templatePath); err != nil {
		_, currentFile, _, _ := runtime.Caller(0)
		templatePath = filepath.Join(filepath.Dir(currentFile), "..", "plantillas", "ticket.html")
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("Error leyendo la plantilla:", err)
		return "", err
	}

	var htmlProcesado bytes.Buffer
	err = tmpl.Execute(&htmlProcesado, datos)
	if err != nil {
		log.Println("Error inyectando datos en la plantilla:", err)
		return "", err
	}

	return htmlProcesado.String(), nil
}

// ConvertirHTMLaImagen envía el HTML a la API y devuelve la URL pública de la imagen
func ConvertirHTMLaImagen(html string) (string, error) {
	apiUserID := os.Getenv("HCTI_USER_ID")
	apiKey := os.Getenv("HCTI_API_KEY")

	// Preparamos los datos para enviarlos por POST
	data := url.Values{}
	data.Set("html", html)

	req, err := http.NewRequest("POST", "https://hcti.io/v1/image", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Autenticación básica requerida por la API
	req.SetBasicAuth(apiUserID, apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Ejecutamos la petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Decodificamos el JSON para extraer la URL de la imagen creada
	var hctiResp HCTIResponse
	json.Unmarshal(body, &hctiResp)

	return hctiResp.URL, nil
}
