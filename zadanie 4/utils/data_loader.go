package utils

import (
	"net/http"
	"fmt"
	"io"
	"strings"
	"os"
)

const dataset1 string = "https://archive.ics.uci.edu/ml/machine-learning-databases/wine-quality/winequality-white.csv"

func LoadData() error {
	response, err := http.Get(dataset1)
	if err != nil {
		return fmt.Errorf("błąd podczas pobierania pliku: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("błąd podczas odczytywania odpowiedzi: %v", err)
		}

		csvContent := string(body)

		csvContent = strings.ReplaceAll(csvContent, ";", ",")

		file, err := os.Create("data.csv")
		if err != nil {
			return fmt.Errorf("błąd podczas tworzenia pliku: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(csvContent)
		if err != nil {
			return fmt.Errorf("błąd podczas zapisywania pliku: %v", err)
		}

		fmt.Println("Plik został pobrany i zmieniony.")
	} else {
		return fmt.Errorf("błąd podczas pobierania pliku: %v", response.StatusCode)
	}

	return nil
}

func Normalize(X [][]float64) [][]float64 {
    normalizedX := make([][]float64, len(X))
    minValues := make([]float64, len(X[0]))
    maxValues := make([]float64, len(X[0]))

    for j := 0; j < len(X[0]); j++ {
        minValues[j] = X[0][j]
        maxValues[j] = X[0][j]
        for i := 1; i < len(X); i++ {
            if X[i][j] < minValues[j] {
                minValues[j] = X[i][j]
            }
            if X[i][j] > maxValues[j] {
                maxValues[j] = X[i][j]
            }
        }
    }

    for i := 0; i < len(X); i++ {
        normalizedX[i] = make([]float64, len(X[i]))
        for j := 0; j < len(X[i]); j++ {
            if maxValues[j] > minValues[j] {
                normalizedX[i][j] = (X[i][j] - minValues[j]) / (maxValues[j] - minValues[j])
            } else {
                normalizedX[i][j] = 0.0
            }
        }
    }

    return normalizedX
}