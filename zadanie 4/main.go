/*
Autorzy:
Jan Wolski s26435
Marcin Topolniak s25672

Polecenie: Użyj drzewa decyzyjnego i SVM do klasyfikcji
- wybierz jeden zbiór danych do klasyfikcaji
- naucz drzewo decyzyjne i SVM klasyfikować dane
- wybierz drugi zbiór dancyh do klasyfikacji
- naucz drzewo decyzyjne i SVM klasyfikować dane
- pokaż metryki związane z jakością klasyfikacji
- podaj przykładową wizualizację danych
- wywołaj klasyfikatory dla przykładowych danych wejściowcych

Instrukcja przygotowania środowiska znajduje się w pliku readme w repozytorium
*/
package main
/*
Plik main.go stanowi punkt wejścia do projektu. 
Pozwala na przeprowadzenie klasyfikacji za pomocą modeli podanych w treści zadania.
*/

import (
	"zad4/models"
)

//wybór numeru setu
const num_set int = 3

//wybrany set danych jest klasyfikowany za pomocą modeli SVM i drzewa decyzyjnego
func main() {
	models.ShowSVM(num_set)
	models.ShowTree(num_set)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
