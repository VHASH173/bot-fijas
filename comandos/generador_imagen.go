package comandos

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

// Estructura de los datos que van al ticket
type DatosTicket struct {
    Partido    string
    Mercado    string
    Pronostico string
    Cuota      string
}

// Genera el HTML inyectando los datos y forzando un fondo blanco/dimensiones
func GenerarHTMLTicket(datos DatosTicket) string {
    return fmt.Sprintf(`
    <html>
    <head>
        <style>
            body { 
                background-color: #1a1a1a; 
                color: white; 
                font-family: Arial, sans-serif; 
                width: 400px; /* Evita el fondo negro gigante */
                height: 250px; 
                margin: 0; 
                padding: 20px; 
                box-sizing: border-box;
                border: 2px solid #00ff88;
                border-radius: 10px;
            }
            .titulo { color: #00ff88; font-weight: bold; text-align: center; font-size: 20px; }
            .detalle { margin-top: 15px; font-size: 16px; }
            .cuota { color: #00ff88; font-weight: bold; font-size: 18px; }
        </style>
    </head>
    <body>
        <div class="titulo">🔥 APUESTA CONFIRMADA 🔥</div>
        <div class="detalle">
            <p><strong>Partido:</strong> %s</p>
            <p><strong>Mercado:</strong> %s</p>
            <p><strong>Pronóstico:</strong> %s</p>
            <p><strong>Cuota:</strong> <span class="cuota">%s</span></p>
        </div>
    </body>
    </html>`, datos.Partido, datos.Mercado, datos.Pronostico, datos.Cuota)
}

// Se conecta a la API de HCTI y devuelve la URL de la imagen
func ConvertirHTMLaImagen(htmlContent string) (string, error) {
    apiID := os.Getenv("HCTI_USER_ID")
    apiKey := os.Getenv("HCTI_API_KEY")

    data := map[string]string{"html": htmlContent}
    jsonData, _ := json.Marshal(data)

    req, err := http.NewRequest("POST", "https://hcti.io/v1/image", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    req.SetBasicAuth(apiID, apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    var result map[string]interface{}
    json.Unmarshal(body, &result)

    if url, ok := result["url"].(string); ok {
        return url, nil
    }
    return "", fmt.Errorf("no se pudo generar la imagen: %s", string(body))
}
