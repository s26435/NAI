package utils

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)
/*
Plik data_set implementuje różne typt zestawów danych oraz operacje na nich
Obsługiwane zestawy:
1. Dataset - dane dotyczące jakości wina
2. DatasetPassenger - dane pasażeró z Titanica
3. DatasetInsulin - dane o insulinach pobierane z API
*/
type Dataset []WineData

type WineData struct {
	FixedAcidity       float64
	VolatileAcidity    float64
	CitricAcid         float64
	ResidualSugar      float64
	Chlorides          float64
	FreeSulfurDioxide  float64
	TotalSulfurDioxide float64
	Density            float64
	PH                 float64
	Sulphates          float64
	Alcohol            float64
	Quality            float64
}


/*
Funkcja print wyświetla zawrtość zestawu danych
- iteruje przez dane i wyświetla szczegóły każdej próbki
- zatrzymuje wyświetlanie po osiągnięciu limitu
*/
func (ds Dataset) Print(optionalArgs ...int) {
	var tr int = -1
	if len(optionalArgs) == 1 {
		tr = optionalArgs[0] - 1
	}
	for i, wine := range ds {
		fmt.Printf("%+v\n", wine)
		if tr != -1 && tr == i {
			break
		}
	}
}

/*
Funkcja Visualize wizualizuje rozkład jakości wina w zestawie danych
- oblicza liczność próbek dla każdej jakości
- tworzy histogram z rozkładem jakości
- zapisuje wykres do pliku
*/
func (ds Dataset) Visualize() {
	// Zliczanie jakości wina
	qualityCounts := make(map[int]int)
	for _, data := range ds {
		qualityCounts[int(data.Quality)]++
	}

	// Przygotowanie danych do wykresu
	bars := make(plotter.Values, len(qualityCounts))
	labels := []string{}
	i := 0
	for quality, count := range qualityCounts {
		bars[i] = float64(count)
		labels = append(labels, fmt.Sprintf("%d", quality))
		i++
	}

	// Tworzenie wykresu
	p, err := plot.New()
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia wykresu: %v\n", err)
		return
	}

	p.Title.Text = "Rozkład jakości wina"
	p.Y.Label.Text = "Liczba próbek"

	hist, err := plotter.NewBarChart(bars, vg.Points(20))
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia histogramu: %v\n", err)
		return
	}

	p.Add(hist)

	// Zapisywanie wykresu do pliku
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "wine_quality_distribution.png"); err != nil {
		fmt.Printf("Błąd podczas zapisywania wykresu: %v\n", err)
		return
	}

	fmt.Println("Wykres zapisany do wine_quality_distribution.png")
}
/* 
Funkcja LoadData pobiera dane o jakości wina z pliku CSV
- pobiera plik CSV z podanego URL
- przetwarza dane CSV, konwertują je na struktury 'WineData'
- obsługuje potencjalne błędy
*/
func (ds *Dataset) LoadData() {
	url := "https://archive.ics.uci.edu/ml/machine-learning-databases/wine-quality/winequality-white.csv"
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Blad podczas wykonywania zapytania GET: %v\n", err)
		return
	}
	defer response.Body.Close()
	file, err := os.Create("data.csv")
	if err != nil {
		fmt.Printf("Blad podczas tworzenia pliku: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("Blad podczas zapisywania do pliku: %v\n", err)
		return
	}

	csvFile, err := os.Open("data.csv")
	if err != nil {
		fmt.Printf("Blad podczas otwierania pliku CSV: %v\n", err)
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Blad podczas odczytu pliku CSV: %v\n", err)
		return
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 12 {
			fmt.Printf("Blad: niewystarczajaca liczba kolumn w wierszu %d\n", i)
			return
		}
		fixedAcidity, _ := strconv.ParseFloat(record[0], 64)
		volatileAcidity, _ := strconv.ParseFloat(record[1], 64)
		citricAcid, _ := strconv.ParseFloat(record[2], 64)
		residualSugar, _ := strconv.ParseFloat(record[3], 64)
		chlorides, _ := strconv.ParseFloat(record[4], 64)
		freeSulfurDioxide, _ := strconv.ParseFloat(record[5], 64)
		totalSulfurDioxide, _ := strconv.ParseFloat(record[6], 64)
		density, _ := strconv.ParseFloat(record[7], 64)
		ph, _ := strconv.ParseFloat(record[8], 64)
		sulphates, _ := strconv.ParseFloat(record[9], 64)
		alcohol, _ := strconv.ParseFloat(record[10], 64)
		quality, _ := strconv.ParseFloat(record[11], 64)

		*ds = append(*ds, WineData{
			FixedAcidity:       fixedAcidity,
			VolatileAcidity:    volatileAcidity,
			CitricAcid:         citricAcid,
			ResidualSugar:      residualSugar,
			Chlorides:          chlorides,
			FreeSulfurDioxide:  freeSulfurDioxide,
			TotalSulfurDioxide: totalSulfurDioxide,
			Density:            density,
			PH:                 ph,
			Sulphates:          sulphates,
			Alcohol:            alcohol,
			Quality:            quality,
		})
	}
}

/*
Funkcja ToXY konwertuje zestaw danych na cechy i etykiety
- ekstraktuje cechy
- konwertuje jakość wina na liczbę załkowitą jako etykietę klasy
*/
func (ds Dataset) ToXY() ([][]float64, []int) {
	X := make([][]float64, len(ds))
	Y := make([]int, len(ds))
	for i, data := range ds {
		X[i] = []float64{
			data.FixedAcidity,
			data.VolatileAcidity,
			data.CitricAcid,
			data.ResidualSugar,
			data.Chlorides,
			data.FreeSulfurDioxide,
			data.TotalSulfurDioxide,
			data.Density,
			data.PH,
			data.Sulphates,
			data.Alcohol,
		}
		Y[i] = int(data.Quality)
	}
	return X, Y
}

/*
Funkcja TrainTestSplit dzieli na dane zestawy treningowe i testowe
- generuje losową permutację indeksów próbek
- dzieli dane na podstawie indeksów
*/
func (ds Dataset) TrainTestSplit(testSize float64) (Dataset, Dataset) {
	rand.Seed(time.Now().UnixNano())
	indices := rand.Perm(len(ds))

	numTest := int(float64(len(ds)) * testSize)
	trainIndices := indices[numTest:]
	testIndices := indices[:numTest]

	trainSet := make(Dataset, len(trainIndices))
	testSet := make(Dataset, len(testIndices))

	for i, idx := range trainIndices {
		trainSet[i] = ds[idx]
	}

	for i, idx := range testIndices {
		testSet[i] = ds[idx]
	}

	return trainSet, testSet
}

type DatasetPasanger []PassengerData

type PassengerData struct {
	Survived int
	Pclass   int
	Sex      int // 0: male, 1: female
	Age      float64
	SibSp    int
	Parch    int
	Fare     float64
}

/*
Funkcja Print wyświetla zawartość danych pasażeró
- wyświetla dane pasażerów
- zatrzymuje wyświetlanie po osiągnięciu limitu
*/
func (ds DatasetPasanger) Print(optionalArgs ...int) {
	var tr int = -1
	if len(optionalArgs) == 1 {
		tr = optionalArgs[0] - 1
	}
	for i, passenger := range ds {
		fmt.Printf("%+v\n", passenger)
		if tr != -1 && tr == i {
			break
		}
	}
}

/*
Funkcja LoadData wczytuje dane pasażerów z pliku CSV
- otwiera plik CSV i odczytuje jego zawartość
- konwertuje dane na strukturę 'PassengerData'
- obsługuje brakujące dane i błędy w formacie pliku
*/
func (ds *DatasetPasanger) LoadData() {
	filename := "tested.csv"
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Blad podczas otwierania pliku CSV: %v\n", err)
		return
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Blad podczas odczytu pliku CSV: %v\n", err)
		return
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 12 {
			fmt.Printf("Blad: niewystarczajaca liczba kolumn w wierszu %d\n", i)
			continue
		}

		survived, _ := strconv.Atoi(record[1])
		pclass, _ := strconv.Atoi(record[2])
		sex := 0
		if record[4] == "female" {
			sex = 1
		}
		age, _ := strconv.ParseFloat(record[5], 64)
		sibSp, _ := strconv.Atoi(record[6])
		parch, _ := strconv.Atoi(record[7])
		fare, _ := strconv.ParseFloat(record[9], 64)
		// fmt.Println(survived, pclass,sex,age,sibSp,parch,fare)
		*ds = append(*ds, PassengerData{
			Survived: survived,
			Pclass:   pclass,
			Sex:      sex,
			Age:      age,
			SibSp:    sibSp,
			Parch:    parch,
			Fare:     fare,
		})
	}
}

/*
Funkcja Visualize wizualizuje rozkład liczy pasażerów, którzy przeżyli lub zmarli
- oblicza liczbę pasażerów dla każdej grupy
- tworzy histogram i zapisuje go
*/
func (ds DatasetPasanger) Visualize() {
	// Zliczanie ocalałych i zmarłych pasażerów
	survivedCount := 0
	notSurvivedCount := 0
	for _, data := range ds {
		if data.Survived == 1 {
			survivedCount++
		} else {
			notSurvivedCount++
		}
	}

	// Przygotowanie danych do wykresu
	values := plotter.Values{float64(survivedCount), float64(notSurvivedCount)}

	// Tworzenie wykresu
	p, err := plot.New()
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia wykresu: %v\n", err)
		return
	}

	p.Title.Text = "Rozkład liczby pasażerów"
	p.Y.Label.Text = "Liczba pasażerów"
	p.NominalX("Ocalałych", "Zmarłych")

	// Histogram
	barChart, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia histogramu: %v\n", err)
		return
	}

	p.Add(barChart)

	// Zapisywanie wykresu do pliku
	err = p.Save(8*vg.Inch, 6*vg.Inch, "titanic_histogram.png")
	if err != nil {
		fmt.Printf("Błąd podczas zapisywania wykresu: %v\n", err)
		return
	}

	fmt.Println("Wykres zapisany do titanic_histogram.png")
}

/*
Funkcja ToXY konwertuje dane pasażerów na cechy i etykiety
- ekstraktuje cechy
- ustawia zmienna 'Survived' jako etykietę klasy
*/
func (ds DatasetPasanger) ToXY() ([][]float64, []int) {
	X := make([][]float64, len(ds))
	Y := make([]int, len(ds))
	for i, data := range ds {
		X[i] = []float64{
			float64(data.Pclass),
			float64(data.Sex),
			data.Age,
			float64(data.SibSp),
			float64(data.Parch),
			data.Fare,
		}
		Y[i] = data.Survived
	}
	return X, Y
}

// Funkcja TrainTestSplit dzieli dane pasażerów na zestawy treningowe i testowe
func (ds DatasetPasanger) TrainTestSplit(testSize float64) (DatasetPasanger, DatasetPasanger) {
	rand.Seed(time.Now().UnixNano())
	indices := rand.Perm(len(ds))

	numTest := int(float64(len(ds)) * testSize)
	trainIndices := indices[numTest:]
	testIndices := indices[:numTest]

	trainSet := make(DatasetPasanger, len(trainIndices))
	testSet := make(DatasetPasanger, len(testIndices))

	for i, idx := range trainIndices {
		trainSet[i] = ds[idx]
	}

	for i, idx := range testIndices {
		testSet[i] = ds[idx]
	}

	return trainSet, testSet
}

type DatasetInsulin []APIData

type APIData struct {
	IsHuman        int
	LineageCount   int
	SequenceLength int
	MolWeight      int
}

// Funkcja Print wyświetla dane dotyczące sekwencji insulin
func (ds DatasetInsulin) Print(optionalArgs ...int) {
	var tr int = -1
	if len(optionalArgs) == 1 {
		tr = optionalArgs[0] - 1
	}
	for i, data := range ds {
		fmt.Printf("%+v\n", data)
		if tr != -1 && tr == i {
			break
		}
	}
}

/*
Funkcja LoadData pobiera dane o insulinach z API UniProt
- wysyła zapytanie HTTP do Api UniProt
- dekoduje odpowiedź JSON i konwertuje dane na strukturę APIData
- obsługuje dekompresje odpowiedzi oraz błędy w połączeniu
*/
func (ds *DatasetInsulin) LoadData() {
	apiURL := "https://rest.uniprot.org/uniprotkb/stream?compressed=false&query=reviewed:true+AND+insulin&fields=organism_id,mass,length,lineage_ids&size=500"

	// Tworzenie zapytania HTTP
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia zapytania HTTP: %v\n", err)
		return
	}

	// Ustawienie nagłówków
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	// Wysłanie zapytania
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Błąd podczas wysyłania zapytania HTTP: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Dekompresja odpowiedzi (gzip lub deflate)
	var reader io.Reader = response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			fmt.Printf("Błąd podczas otwierania gzip: %v\n", err)
			return
		}
		defer gzipReader.Close()
		reader = gzipReader
	} else if response.Header.Get("Content-Encoding") == "deflate" {
		fmt.Println("Obsługa formatu deflate nie jest jeszcze zaimplementowana.")
		return
	}

	// Odczytanie danych do bufora
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Błąd podczas odczytu odpowiedzi: %v\n", err)
		return
	}

	// Dekodowanie JSON
	var apiResponse struct {
		Results []struct {
			Organism struct {
				TaxonID int `json:"taxonId"`
			} `json:"organism"`
			Sequence struct {
				Length    int `json:"length"`
				MolWeight int `json:"molWeight"`
			} `json:"sequence"`
			Lineages []struct {
				TaxonID int `json:"taxonId"`
			} `json:"lineages"`
		} `json:"results"`
	}

	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		fmt.Printf("Błąd podczas parsowania JSON: %v\n", err)
		return
	}

	// Przetwarzanie wyników
	for _, result := range apiResponse.Results {
		isHuman := 0
		if result.Organism.TaxonID == 9606 {
			isHuman = 1
		}

		lineageCount := len(result.Lineages)

		*ds = append(*ds, APIData{
			IsHuman:        isHuman,
			LineageCount:   lineageCount,
			SequenceLength: result.Sequence.Length,
			MolWeight:      result.Sequence.MolWeight,
		})
	}
}

/*
Funkcja Visualize wizualizuje rozkłąd długości sekwencji insulin
- oblicza liczność próbek dla każdej długości sekwencji
- tworzy histogram i zapisuje go
*/
func (ds DatasetInsulin) Visualize() {
	sequenceLengths := make(map[int]int)
	for _, data := range ds {
		sequenceLengths[data.SequenceLength]++
	}

	// Przygotowanie danych do wykresu
	bars := make(plotter.Values, len(sequenceLengths))
	labels := []string{}
	i := 0
	for length, count := range sequenceLengths {
		bars[i] = float64(count)
		labels = append(labels, fmt.Sprintf("%d", length))
		i++
	}

	// Tworzenie wykresu
	p, err := plot.New()
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia wykresu: %v\n", err)
		return
	}

	p.Title.Text = "Rozkład długości sekwencji insulin"
	p.Y.Label.Text = "Liczba próbek"

	hist, err := plotter.NewBarChart(bars, vg.Points(20))
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia histogramu: %v\n", err)
		return
	}

	p.Add(hist)

	// Zapisywanie wykresu do pliku
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "insulin_sequence_length_distribution.png"); err != nil {
		fmt.Printf("Błąd podczas zapisywania wykresu: %v\n", err)
		return
	}

	fmt.Println("Wykres zapisany do insulin_sequence_length_distribution.png")
}

/*
Funkcja ToXY konweruje dane o insulinach na cechy i etykiety
- ekstrahuje cechy
- ustawia zmienną IsHuman jako etykietę klasy
*/
func (ds DatasetInsulin) ToXY() ([][]float64, []int) {
	X := make([][]float64, len(ds))
	Y := make([]int, len(ds))
	for i, data := range ds {
		X[i] = []float64{
			float64(data.LineageCount),
			float64(data.SequenceLength),
			float64(data.MolWeight),
		}
		Y[i] = data.IsHuman
	}
	return X, Y
}

// Funkcja TrainTestSplit dzieli dane o insulinach na zestawy treningowe i testowe
func (ds DatasetInsulin) TrainTestSplit(testSize float64) (DatasetInsulin, DatasetInsulin) {
	rand.Seed(time.Now().UnixNano())
	indices := rand.Perm(len(ds))

	numTest := int(float64(len(ds)) * testSize)
	trainIndices := indices[numTest:]
	testIndices := indices[:numTest]

	trainSet := make(DatasetInsulin, len(trainIndices))
	testSet := make(DatasetInsulin, len(testIndices))

	for i, idx := range trainIndices {
		trainSet[i] = ds[idx]
	}

	for i, idx := range testIndices {
		testSet[i] = ds[idx]
	}

	return trainSet, testSet
}
