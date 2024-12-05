package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
)

/*Funkcja ShowDiagram generuje plik HTML z diagramem Mermaid i otwiera go w przeglądarce
- wczytuje zawartość pliku .mermaid
- generuje plik html z osadzonym idagramem
- otwiera plik html w domyślnej przeglądarce
*/
func ShowDiagram(filePath string) {
	// Wczytanie plik .mermaid
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Błąd podczas odczytu pliku: %v", err)
	}

	// Generowanie zawartość HTML
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Mermaid Diagram</title>
		<script type="module" src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
		<script>
			mermaid.initialize({ startOnLoad: true });
		</script>
	</head>
	<body>
		<div class="mermaid">
			%s
		</div>
	</body>
	</html>
	`, string(data))

	// Zapisanie HTML do pliku
	outputFile := "diagram.html"
	err = ioutil.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("Błąd podczas zapisu pliku HTML: %v", err)
	}

	// Otwarcie plik w przeglądarce
	openBrowser(outputFile)
	fmt.Println("Diagram został wyświetlony w przeglądarce.")
}

/* Funkcja openBrowser otwiera podany plik/URL w przeglądarce
- sprawdza system operacyjny
- używa odpowiedniego polecenia do otwarcia pliku w przeglądarce
*/
func openBrowser(fileName string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", fileName).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fileName).Start()
	case "darwin":
		err = exec.Command("open", fileName).Start()
	default:
		log.Fatalf("Nieobsługiwany system operacyjny")
	}

	if err != nil {
		log.Fatalf("Błąd podczas otwierania przeglądarki: %v", err)
	}
}

// Funkcja Must wywołuje panic w przypadku błędu
func Must(err error){
	if err != nil {
		panic(err)
	}
}
