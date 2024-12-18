package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"log"

	"github.com/patrikeh/go-deep/training"
)

type FashionDataset struct {
	Images [][]float64
	Labels [][]float64
	Rows   int
	Cols   int
}

// FashionLabels zawiera mapowanie etykiet liczbowych na ich nazwy tekstowe.
var FashionLabels = map[int]string{
	0: "T-shirt/top",
	1: "Trouser",
	2: "Pullover",
	3: "Dress",
	4: "Coat",
	5: "Sandal",
	6: "Shirt",
	7: "Sneaker",
	8: "Bag",
	9: "Ankle boot",
}

var Labels = []string {"T-shirt/top", "Trouser", "Pullover", "Dress", "Coat", "Sandal", "Shirt", "Sneaker", "Bag", "Ankle boot"}

// LoadData wczytuje dane obrazów i etykiet z plików w formacie ubyte.
// Wczytuje dane z plików "images-ubyte" i "labels-ubyte".
func (fd *FashionDataset) LoadData() {
	const imagePath, labelPath string = "images-ubyte", "labels-ubyte"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Plik %s nie istnieje", imagePath)
	}
	if _, err := os.Stat(labelPath); os.IsNotExist(err) {
		log.Fatalf("Plik %s nie istnieje", labelPath)
	}

	imageFile, err := os.Open(imagePath)
	Must(err)
	defer imageFile.Close()

	var magic, numImages, numRows, numCols int32
	err = binary.Read(imageFile, binary.BigEndian, &magic)
	Must(err)
	if magic != 2051 {
		Must(fmt.Errorf("nieprawidłowy magic number dla obrazów: %d", magic))
	}
	
	binary.Read(imageFile, binary.BigEndian, &numImages)
	binary.Read(imageFile, binary.BigEndian, &numRows)
	binary.Read(imageFile, binary.BigEndian, &numCols)

	fd.Rows = int(numRows)
	fd.Cols = int(numCols)
	fd.Images = make([][]float64, numImages)

	for i := int32(0); i < numImages; i++ {
		img := make([]float64, numRows*numCols)
		pixelData := make([]byte, numRows*numCols)
		err := binary.Read(imageFile, binary.BigEndian, pixelData)
		Must(err)
		for j := range pixelData {
			img[j] = float64(pixelData[j]) / 255.0
		}
		fd.Images[i] = img
	}

	labelFile, err := os.Open(labelPath)
	Must(err)
	defer labelFile.Close()

	err = binary.Read(labelFile, binary.BigEndian, &magic)
	Must(err)
	if magic != 2049 {
		Must(fmt.Errorf("nieprawidłowy magic number dla labeli: %d", magic))
	}

	var numLabels int32
	binary.Read(labelFile, binary.BigEndian, &numLabels)
	if numLabels != numImages {
		Must(fmt.Errorf("niezgodność liczby obrazów i labeli: %d vs %d", numImages, numLabels))
	}

	labelData := make([]byte, numLabels)
	Must(binary.Read(labelFile, binary.BigEndian, labelData))

	fd.Labels = make([][]float64, numLabels)
	for i, label := range labelData {
		fd.Labels[i] = GetHotOneFloat64(label, 10)
	}
}

// ToExamples konwertuje zbior danych na format wymagany przez bibliotekę treningową.
// Zwraca tablicę training.Example zawierającą dane wejściowe i etykiety.
func (fd *FashionDataset) ToExamples() []training.Example {
	examples := make([]training.Example, len(fd.Images))

	for i, image := range fd.Images {
		examples[i] = training.Example{
			Input:    image,
			Response: fd.Labels[i], // One-hot etykieta jako output
		}
	}
	return examples
}

// Print wypisuje obraz i jego etykietę na konsoli.
// Można podać opcjonalny indeks obrazu do wydrukowania.
func (fd *FashionDataset) Print(optionalArgs ...int) {
	index := 0
	if len(optionalArgs) > 0 {
		index = optionalArgs[0]
	}
	if index < 0 || index >= len(fd.Images) {
		fmt.Println("Nieprawidłowy indeks")
		return
	}

	ones := fd.Labels[index]
	labelIndex := -1
	for i, val := range ones {
		if val == 1.0 {
			labelIndex = i
			break
		}
	}
	labelName, exists := FashionLabels[labelIndex]
	if !exists {
		labelName = "Unknown"
	}

	fmt.Printf("Etykieta: %d (%s)\n", labelIndex, labelName)
	fmt.Println("Obraz: ")
	for r := 0; r < fd.Rows; r++ {
		for c := 0; c < fd.Cols; c++ {
			fmt.Printf("%0.2f ", fd.Images[index][r*fd.Cols+c])
		}
		fmt.Println()
	}
}

// ToXY zwraca dane wejściowe (X) i etykiety (Y) jako dwie tablice.
func (fd *FashionDataset) ToXY() ([][]float64, [][]float64) {
	return fd.Images, fd.Labels
}

// Normalize normalizuje wartości pikseli obrazów do zakresu [0, 1].
// Zwraca błąd, jeśli zbior danych jest pusty.
func (fd FashionDataset) Normalize() error {
	if len(fd.Images) == 0 {
		return fmt.Errorf("brak danych do normalizacji")
	}

	for i := range fd.Images {
		for j := range fd.Images[i] {
			fd.Images[i][j] = fd.Images[i][j] / math.Max(1.0, fd.Images[i][j])
		}
	}
	return nil
}

// TrainTestSplit dzieli zbior danych na zestaw treningowy i testowy.
// testSize - ułamek danych przeznaczony na zestaw testowy.
func (fd *FashionDataset) TrainTestSplit(testSize float64) (FashionDataset, FashionDataset) {
	indices := rand.Perm(len(fd.Images))
	numTest := int(float64(len(fd.Images)) * testSize)

	testIndices := indices[:numTest]
	trainIndices := indices[numTest:]

	trainSet := &FashionDataset{
		Images: make([][]float64, len(trainIndices)),
		Labels: make([][]float64, len(trainIndices)),
		Rows:   fd.Rows,
		Cols:   fd.Cols,
	}
	testSet := &FashionDataset{
		Images: make([][]float64, len(testIndices)),
		Labels: make([][]float64, len(testIndices)),
		Rows:   fd.Rows,
		Cols:   fd.Cols,
	}

	for i, idx := range trainIndices {
		trainSet.Images[i] = fd.Images[idx]
		trainSet.Labels[i] = fd.Labels[idx]
	}

	for i, idx := range testIndices {
		testSet.Images[i] = fd.Images[idx]
		testSet.Labels[i] = fd.Labels[idx]
	}

	return *trainSet, *testSet
}