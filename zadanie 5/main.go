/*
Autorzy:
Jan Wolski s26435
Marcin Topolniak s25672

Polecenie:
- Wykorzystać jeden z zbiorów danych z poprzednich ćwiczeń i nauczyć sieć neuronową.
- Porównać skuteczność obu podejść. Dodać logi/print screen do repozytorium.
- Nauczyć sieć rozpoznawać zwierzęta, np. z zbioru CIFAR10.
- Nauczyć sieć rozpoznawać ubrania, np. GitHub - zalandoresearch/fashion-mnist: A MNIST-like fashion product database.
- Zaskocz mnie. Zaproponuj własny przypadek użycia sieci neuronowych do problemu klasyfikacji.
- Dla jednego z punktu narysuj confusion matrix.
- Do jednego z punktów użyj dwóch rozmiarów sieci neuronowych. Porównaj wyniki. 

Instrukcja przygotowania środowiska znajduje się w pliku readme w repozytorium
*/

package main

import (
	"venv/models"
)

/*
Plik main.go stanowi punkt wejścia do projektu. 
Pozwala na przeprowadzenie klasyfikacji za pomocą modeli podanych w treści zadania.
*/

/*
w package models są przygotowane funkcje demonstracyjne rozwiązań zadań.
Funkcje:
- models.ShowInsuliNN()
- models.ShowFashionNN()
- models.ShowCifarNN()
- models.ShowAlzhaimer()

Wszystkie zawierają Confusion matrix (w poleceniu wymagana był tylko jedna)
Porównanie dwóch różnych wielkości sieci zawiera ShowFashionNN()
*/

func main(){
	models.ShowAlzhaimer()
}