package utils

import (
	"fmt"
	"math"
	"os"

	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)


// EvaluateModelInsulin ocenia dokładność modelu predykcyjnego na podstawie macierzy pomyłek.
// data - przykłady danych testowych typu training.Examples
// network - sieć neuronowa typu deep.Neural, która dokonuje predykcji.
func EvaluateModelInsulin(data training.Examples, network *deep.Neural) {
	truePositive := 0
	trueNegative := 0
	falsePositive := 0
	falseNegative := 0

	for _, example := range data {
		prediction := network.Predict(example.Input)
		predictedLabel := 0
		if prediction[0] >= 0.5 { 
			predictedLabel = 1
		}

		trueLabel := int(example.Response[0])
		if predictedLabel == 1 && trueLabel == 1 {
			truePositive++
		} else if predictedLabel == 0 && trueLabel == 0 {
			trueNegative++
		} else if predictedLabel == 1 && trueLabel == 0 {
			falsePositive++
		} else if predictedLabel == 0 && trueLabel == 1 {
			falseNegative++
		}
	}

	total := truePositive + trueNegative + falsePositive + falseNegative
	accuracy := float64(truePositive+trueNegative) / float64(total)

	fmt.Println("Confusion Matrix:")
	fmt.Printf("TP: %d, FP: %d\n", truePositive, falsePositive)
	fmt.Printf("FN: %d, TN: %d\n", falseNegative, trueNegative)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)
}

// GetHotOne konwertuje liczbę etykiety na kodowanie one-hot.
// label - wartość etykiety.
// numClasses - liczba klas.
// Zwraca tablicę int z zakodowaną wartością etykiety.
func GetHotOne(label int, numClasses int) []int {
	res := make([]int, numClasses)
	res[int(label)] = 1
	return res
}

// convertToOneHot konwertuje etykietę liczbową na format one-hot encoding.
// label - wartość etykiety (0-9)
// numClasses - liczba klas
// Zwraca tablicę float64 z zakodowaną wartością etykiety.
func GetHotOneFloat64(label byte, numClasses int) []float64 {
	oneHot := make([]float64, numClasses)
	oneHot[int(label)] = 1.0
	return oneHot
}

// GetFloatArr konwertuje tablicę liczb całkowitych na tablicę liczb zmiennoprzecinkowych.
// arr - wejściowa tablica liczb całkowitych.
// Zwraca tablicę liczb zmiennoprzecinkowych odpowiadającą wejściowej tablicy.
func GetFloatArr(arr []int) []float64 {
    result := make([]float64, len(arr))
    for i, val := range arr {
        result[i] = float64(val)
    }
    return result
}

// EvaluateModelCifar ocenia dokładność modelu na zbiorze danych testowych CIFAR.
// Funkcja oblicza dokładność klasyfikacji modelu, generuje macierz błędów (confusion matrix) i wyświetla wyniki na standardowym wyjściu. 
// testData []training.Example: Zbiór przykładów testowych zawierających dane wejściowe oraz oczekiwane odpowiedzi.
// model *deep.Neural: Wskaźnik na model sieci neuronowej, który ma zostać oceniony.
// numClasses int: Liczba klas występujących w zadaniu klasyfikacji.
// classLabels []string: Lista etykiet klas, które będą użyte do wyświetlenia macierzy błędów.
func EvaluateModelCifar(testData []training.Example, model *deep.Neural, numClasses int, classLabels []string) {
	correct := 0
	confusionMatrix := make([][]int, numClasses)
	for i := range confusionMatrix {
		confusionMatrix[i] = make([]int, numClasses)
	}

	for _, example := range testData {
		prediction := model.Predict(example.Input)
		predictedClass := ArgMax(prediction)
		actualClass := ArgMax(example.Response)

		confusionMatrix[actualClass][predictedClass]++

		if predictedClass == actualClass {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(testData)) * 100
	fmt.Printf("Test Accuracy: %.2f%%\n", accuracy)

	fmt.Println("Confusion Matrix:")
	fmt.Printf("%12s", "Predicted")
	for i := 0; i < numClasses; i++ {
		fmt.Printf("%12s", classLabels[i])
	}
	fmt.Println()

	for i, row := range confusionMatrix {
		fmt.Printf("%-12s", classLabels[i]) 
		for _, count := range row {
			fmt.Printf("%12d", count)
		}
		fmt.Println()
	}
}

// EvaluateAlzhaimerModel ocenia dokładność modelu predykcyjnego na podstawie dokładności modelu.
// data - przykłady danych testowych typu training.Examples
// network - sieć neuronowa typu deep.Neural, która dokonuje predykcji.
func EvaluateAlzhaimerModel(network *deep.Neural, testData training.Examples) {
	var correct int
	var threshold = 0.2 

	for _, example := range testData {
		output := network.Predict(example.Input)
		predicted := output[0]
		actual := example.Response[0]

		if math.Abs(predicted-actual) < threshold {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(testData)) * 100
	fmt.Printf("Dokładność modelu: %.2f%%\n", accuracy)
}




// ArgMax zwraca indeks największej wartości w tablicy liczb zmiennoprzecinkowych.
// values - tablica liczb zmiennoprzecinkowych.
// Zwraca indeks maksymalnej wartości.
func ArgMax(values []float64) int {
	maxIndex := 0
	maxValue := math.Inf(-1)

	for i, value := range values {
		if value > maxValue {
			maxValue = value
			maxIndex = i
		}
	}

	return maxIndex
}

// CheckIfOne sprawdza, czy wszystkie wartości w dwuwymiarowej tablicy mieszczą się w przedziale [0,1].
// img - dwuwymiarowa tablica liczb zmiennoprzecinkowych.
// Zwraca true, jeśli wszystkie wartości są w przedziale [0,1], w przeciwnym razie false.
func CheckIfOne(img [][]float64)bool{
	for _, x := range img{
		for _, j := range x{
			if !(j >= 0 && j<=1){
				return false
			}
		}
	}
	return true
}

// getNewInsulin tworzy nowy obiekt DatasetInsulin.
// Zwraca nowy obiekt DatasetInsulin.
func getNewInsulin()DatasetInsulin{ 
	return DatasetInsulin{}
}

// getNewCifar tworzy nowy obiekt CifarDataSet.
// Zwraca nowy obiekt CifarDataSet.
func getNewCifar()CifarDataSet{
	return CifarDataSet{}
}

// getNewCifar tworzy nowy obiekt FashionDataset.
// Zwraca nowy obiekt FashionDataset.
func getNewFashion()FashionDataset{ //ignore
	return FashionDataset{}
}

// GetDataset zwraca nowy zestaw danych na podstawie podanej nazwy.
// datasetName - nazwa zestawu danych (np. "uniprot" lub "cifar").
// Zwraca obiekt implementujący DatasetInterface.
// func GetDataset(datasetName string)DatasetInterface{
// 	switch(datasetName){
// 	// case "uniprot":
// 	// 	return getNewInsulin()
// 	case "cifar":
// 		return getNewCifar()
// 	// case "fashion":
// 	// 	return getNewFashion()
// 	default:
// 		panic("wrong dataset name")
// 	}
// }

// Must sprawdza, czy wystąpił błąd. W przypadku błędu program zostaje zakończony.
// err - obiekt błędu do sprawdzenia.
func Must(err error) {
	if err != nil {
		fmt.Printf("Wystąpił błąd: %v\n", err)
		os.Exit(1)
	}
}
