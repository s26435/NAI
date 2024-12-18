package utils

import (
	_ "encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/patrikeh/go-deep/training"
)

type CifarRecord struct {
	Label []int      
	Image []float64 
}

var CifarLabels = []string{"0","1", "2", "3", "4", "5", "6", "7", "8", "9"}

type CifarDataSet []CifarRecord

// Print wypisuje etykiety i rozmiar obrazów w zestawie danych. 
// Przyjmuje opcjonalny argument, który określa liczbę wierszy do wypisania.
func (ds CifarDataSet) Print(optionalArgs ...int) {
	var tr int = -1
	if len(optionalArgs) == 1 {
		tr = optionalArgs[0] - 1
	}
	for i, data := range ds {
		fmt.Printf("%+v: %d\n", data.Label, len(data.Image))
		if tr != -1 && tr == i {
			break
		}
	}
}

// LoadData wczytuje dane CIFAR-10 z plików binarnych data_batch_1.bin - data_batch_5.bin.
// Funkcja odczytuje etykiety oraz dane obrazu i zapisuje je w formacie znormalizowanym.
func (c *CifarDataSet) LoadData() {
	for batch := 1; batch <= 5; batch++ {
		filepath := fmt.Sprintf("data_batch_%d.bin", batch)
		file, err := os.Open(filepath)
		if err != nil {
			panic(fmt.Sprintf("Error opening file %s: %v", filepath, err))
		}
		defer file.Close()

		for {
			label := make([]byte, 1)
			_, err := file.Read(label)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(fmt.Sprintf("Error reading label from %s: %v", filepath, err))
			}

			image := make([]byte, 3072)
			_, err = file.Read(image)
			if err != nil {
				panic(fmt.Sprintf("Error reading image from %s: %v", filepath, err))
			}
			d := CifarRecord{
				Label: GetHotOne(int(label[0]), 10),
				Image: normalizeImage(image),
			}
			*c = append(*c, d)
		}
	}
}

// normalizeImage normalizuje wartości pikseli obrazu z zakresu [0, 255] na [0.0, 1.0].
func normalizeImage(image []byte) []float64 {
	normalized := make([]float64, len(image))
	for i, pixel := range image {
		normalized[i] = float64(pixel) / 255.0
	}
	return normalized
}

// Normalize ponownie normalizuje obrazy w zestawie danych. 
// Funkcja zwraca błąd, jeśli zestaw danych jest pusty.
func (c CifarDataSet) Normalize() error {
	if len(c) == 0 {
		return errors.New("zestaw danych jest pusty, wczytaj dane przed normalizacją")
	}
	for i := range c {
		c[i].Image = normalizeImage(fromFloat64ToByte(c[i].Image))
	}
	return nil
}

// fromFloat64ToByte konwertuje piksele w formacie float64 na format byte.
func fromFloat64ToByte(image []float64) []byte {
	result := make([]byte, len(image))
	for i, value := range image {
		result[i] = byte(value * 255.0)
	}
	return result
}

// ToXY konwertuje zestaw danych na dwie osobne tablice: 
// X - tablica obrazów, Y - tablica etykiet.
func (c CifarDataSet) ToXY() ([][]float64, interface{}) {
	var X [][]float64
	var Y [][]int
	for _, record := range c {
		X = append(X, record.Image)
		Y = append(Y, record.Label)
	}
	return X, Y
}

// ToExamples konwertuje zestaw danych na format zgodny z biblioteką go-deep.
func (c CifarDataSet) ToExamples() []training.Example {
	examples := training.Examples{}
	for _, record := range c {
		examples = append(examples, training.Example{
			Input:    record.Image,
			Response: GetFloatArr(record.Label),
		})
	}
	return examples
}

// TrainTestSplit dzieli zestaw danych na zbiór treningowy i testowy na podstawie podanej proporcji testSize.
func (ds CifarDataSet) TrainTestSplit(testSize float64) (CifarDataSet, CifarDataSet) {
	indices := rand.Perm(len(ds))

	numTest := int(float64(len(ds)) * testSize)
	trainIndices := indices[numTest:]
	testIndices := indices[:numTest]

	trainSet := make(CifarDataSet, len(trainIndices))
	testSet := make(CifarDataSet, len(testIndices))

	for i, idx := range trainIndices {
		trainSet[i] = ds[idx]
	}

	for i, idx := range testIndices {
		testSet[i] = ds[idx]
	}

	return trainSet, testSet
}

