package utils

import (
	"fmt"
	"net/http"
	"math/rand"
	"io"
	"compress/gzip"
	"encoding/json"
	"github.com/patrikeh/go-deep/training"
)

type DatasetInterface interface{
	Print(...int)
	LoadData()
	ToXY()([][]float64, interface{})
	ToExamples()[]training.Example
	Normalize()error
}

type DatasetInsulin []APIData

type APIData struct {
	IsHuman        int
	LineageCount   int
	SequenceLength int
	MolWeight      int
}

// Print wyświetla zawartość zbioru danych. Opcjonalnie może wyświetlić określoną liczbę rekordów.
// Jeśli nie podano argumentu, wyświetlana jest cała zawartość. Jeżeli podano argument `n`, wyświetlane jest pierwsze `n` rekordów.
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

// LoadData pobiera dane o insulinie z publicznego API i zapisuje je w zbiorze.
// Dane są pobierane w formacie JSON i przekształcane do struktury DatasetInsulin.
// Obsługuje kompresję gzip.
func (ds *DatasetInsulin) LoadData() {
	apiURL := "https://rest.uniprot.org/uniprotkb/stream?compressed=false&query=reviewed:true+AND+insulin&fields=organism_id,mass,length,lineage_ids&size=500"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia zapytania HTTP: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Błąd podczas wysyłania zapytania HTTP: %v\n", err)
		return
	}
	defer response.Body.Close()

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

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Błąd podczas odczytu odpowiedzi: %v\n", err)
		return
	}

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

// ToXY konwertuje zbiór danych na macierz cech X i wektor odpowiedzi Y.
// X zawiera cechy: LineageCount, SequenceLength, MolWeight.
// Y zawiera odpowiedzi: 1 dla człowieka, 0 dla innych organizmów.
func (ds DatasetInsulin) ToXY() ([][]float64, interface{}) {
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

// TrainTestSplit dzieli zbiór danych na zestaw treningowy i testowy.
// Parametr testSize określa proporcję danych, które zostaną przeznaczone na zbiór testowy.
func (ds DatasetInsulin) TrainTestSplit(testSize float64) (DatasetInsulin, DatasetInsulin) {
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

// ToExamples konwertuje zbiór danych na format zgodny z biblioteką `go-deep`.
// Każdy rekord jest przekształcany do struktury training.Example.
func (ds DatasetInsulin) ToExamples() []training.Example {
	examples := make([]training.Example, len(ds))

	for i, data := range ds {
		examples[i] = training.Example{
			Input: []float64{float64(data.LineageCount), float64(data.SequenceLength), float64(data.MolWeight)},
			Response: []float64{float64(data.IsHuman)},
		}
	}

	return examples
}

// Normalize normalizuje wartości cech w zbiorze danych do zakresu 0-1.
// Normalizacja jest wykonywana dla LineageCount, SequenceLength i MolWeight.
func (ds DatasetInsulin) Normalize() error {
	if len(ds) == 0 {
		return fmt.Errorf("dataset is empty")
	}

	var minLineage, maxLineage int = (ds)[0].LineageCount, (ds)[0].LineageCount
	var minSeqLength, maxSeqLength int = (ds)[0].SequenceLength, (ds)[0].SequenceLength
	var minMolWeight, maxMolWeight int = (ds)[0].MolWeight, (ds)[0].MolWeight

	for _, data := range ds {
		if data.LineageCount < minLineage {
			minLineage = data.LineageCount
		}
		if data.LineageCount > maxLineage {
			maxLineage = data.LineageCount
		}

		if data.SequenceLength < minSeqLength {
			minSeqLength = data.SequenceLength
		}
		if data.SequenceLength > maxSeqLength {
			maxSeqLength = data.SequenceLength
		}

		if data.MolWeight < minMolWeight {
			minMolWeight = data.MolWeight
		}
		if data.MolWeight > maxMolWeight {
			maxMolWeight = data.MolWeight
		}
	}

	for i := range ds {
		(ds)[i].LineageCount = normalize((ds)[i].LineageCount, minLineage, maxLineage)
		(ds)[i].SequenceLength = normalize((ds)[i].SequenceLength, minSeqLength, maxSeqLength)
		(ds)[i].MolWeight = normalize((ds)[i].MolWeight, minMolWeight, maxMolWeight)
	}
	return nil
}

// MustNormalize wykonuje normalizację zbioru danych i panikuje w przypadku błędu.
func (ds *DatasetInsulin)MustNormalize(){
	err := (*ds).Normalize()
	if err != nil{
		panic(err)
	}
}

func normalize(value, min, max int) int {
	if max == min {
		return 0
	}
	return (value - min) / (max - min)
}