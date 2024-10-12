# NAI - Rozwiązania zadań
Autorzy rozwiązań:
1. Jan Wolski s26435
2. Marcin Topolniak s25672

## Zadanie 1
* Zaimplementować grę dwuosobową, turowa o sumowie zerowej (w naszym przypadku gra Nim)
* Zaimplementować sztuczną inteligencję, który będzie grać w zaproponowaną grę
* Na początku kodu dodać link do zasad; autorzy; instrukcja przygotowania środowiska
* Dodać dokumentację do kodu źródłowego (np. Python -> docstring; Java -> Jdoc; Kotlina -> Kdoc)
* Dodaj zrzu ekranu z przykładowej rozgrywki
* Upewnij się, że nikt z grupy nie wybrał tej samej gry i technologii do rozwiązania
* Prześlij link / dodaj prowadzącego do repozytorium pczapiewski@pjwstk.edu.pl


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

2. **Ustawienie zmiennej środowiskowej GOPATH**:
   - Ustaw zmienną `GOPATH` i dodaj Go do `PATH`:
     ```bash
     echo "export GOPATH=$HOME/go" >> ~/.bashrc
     echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
     source ~/.bashrc
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
   go run main.go
   ```

### Przygotowanie środowiska dla systemu Windows

#### Wymagania wstępne

1. **Zainstaluj Go**:
   - Pobierz instalator Go dla Windows z [oficjalnej strony Go](https://golang.org/dl/).
   - Uruchom instalator i postępuj zgodnie z instrukcjami.

2. **Ustawienie zmiennej środowiskowej GOPATH**:
   - Otwórz Panel sterowania > System > Zaawansowane ustawienia systemu > Zmienne środowiskowe.
   - Dodaj nową zmienną o nazwie `GOPATH` i ustaw wartość na `C:\Users\<TwojaNazwaUżytkownika>\go`.
   - Dodaj `C:\Go\bin` oraz `C:\Users\<TwojaNazwaUżytkownika>\go\bin` do zmiennej `PATH`.

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
   go run main.go
   ```

### Dodatkowe informacje

- Jeśli będziesz potrzebować zainstalować dodatkowe biblioteki, użyj polecenia:
  ```bash
  go get <nazwa_biblioteki>
  ```
