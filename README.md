# NAI - Rozwiązania zadań
Autorzy rozwiązań:
1. Jan Wolski s26435
2. Marcin Topolniak s25672

## Zadanie 1
* Zaimplementować grę dwuosobową, turowa o sumowie zerowej (w naszym przypadku gra Nim)
* Zaimplementować sztuczną inteligencję, który będzie grać w zaproponowaną grę
* Na początku kodu dodać link do zasad; autorzy; instrukcja przygotowania środowiska
* Dodać dokumentację do kodu źródłowego (<a href="https://go.dev/blog/godoc">godoc</a>)
* Dodaj zrzu ekranu z przykładowej rozgrywki
* Upewnij się, że nikt z grupy nie wybrał tej samej gry i technologii do rozwiązania
* Prześlij link / dodaj prowadzącego do repozytorium pczapiewski@pjwstk.edu.pl

## Zadanie 2
* Zaimplementuj system z użyciem logiki rozmytej
* Dodaj dokumentację do kodu źródłowego
* Na początku kodu dodaj opis problemu; wymień autorów rozwiązania; dodaj instrukcję przygotowania środowisk
* Opracowany system musi posiadać przynajmniej 3 wejścia i minimum 1 wyjście
* Wyjście musi być opisane za pomocą przynajmnniej 3 membershio function. Zasada nie obowiązuje, jeżeli system ma więcej wyjść niż 1
* Rozwiązanie umieść w repozytorium z poprzednich ćwiczeń
* Do rozwiązania proszę dołączyć zrzut ekranu/log z przykładowymi 2 wywołaniami systemu
* BONUS: zademonstruj użycie swojego algorytmu do rozwiązania w czasie rzeczywistym
* Problem i technologia muszą być unikatowe w obszarze grupy

## Zadanie 3
* Polecenie: Zaimplementuj silnik rekomandacji filmów/seriali.
* Przestudiuj materiał	A Comparative Study of Clustering Algorithms | by ishika chatterjee | Analytics Vidhya | Medium
* Rozważ uzupełnienie ankiety (samodzielnie)
* Zbuduj silnik rekomendacji filmów i/lub seriali.
* Zaproponuj 5 filmów interesujących dla wybranego użytkownika, których nie oglądał.
* Zaproponouj 5 film, których dany użytkownik nie powinnien oglądać (antyrekomendacje).

### Przygotowanie środowiska dla systemu Linux

#### Wymagania wstępne

1. **Zainstaluj Go**:
   - Otwórz terminal i zainstaluj Go. Możesz to zrobić, korzystając z menedżera pakietów (np. `apt` dla Ubuntu):
     ```bash
     sudo apt update
     sudo apt install golang-go
     ```
   - Po zainstalowaniu sprawdź wersję Go:
     ```bash
     go version
     ```

#### Krok po kroku

1. **Utwórz katalog projektu**:
   ```bash
   mkdir -p ~/go/src/nim-game
   cd ~/go/src/nim-game
   ```

2. **Utwórz plik z kodem**:
   ```bash
   touch main.go
   ```
   - Otwórz plik w edytorze tekstu i wklej kod gry.

3. **Uruchomienie projektu**:
   ```bash
   go mod init mod
   go mod tidy
   go run main.go
   ```

### Przygotowanie środowiska dla systemu Windows

#### Wymagania wstępne

1. **Zainstaluj Go**:
   - Pobierz instalator Go dla Windows z [oficjalnej strony Go](https://golang.org/dl/).
   - Uruchom instalator i postępuj zgodnie z instrukcjami.

#### Krok po kroku

1. **Utwórz katalog projektu**:
   - Otwórz Wiersz polecenia (cmd) i utwórz nowy katalog:
     ```cmd
     mkdir C:\Users\<TwojaNazwaUżytkownika>\go\src\nim-game
     cd C:\Users\<TwojaNazwaUżytkownika>\go\src\nim-game
     ```

2. **Utwórz plik z kodem**:
   - Utwórz plik `main.go`:
     ```cmd
     echo > main.go
     ```
   - Otwórz plik w edytorze tekstu (np. Notepad) i wklej kod gry.

3. **Uruchomienie projektu**:
   ```cmd
   go mod init mod
   go mod tidy
   go run main.go
   ```
