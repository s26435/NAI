package models

import (
	"fmt"
	"math"
	"os"

	"zad4/utils"
)

const file_name string = "tree.mermaid"

// Struktura reprezentująca węzeł w drzewie decyzyjnym
type Node struct {
	Feature   int
	Threshold float64
	Left      *Node
	Right     *Node
	Label     *int
}

// Struktura dla drzewa decyzyjnego
type DecisionTree struct {
	MaxDepth int
	Tree     *Node
}

func getDataset(dataSetNum int)([][]float64, []int, [][]float64, []int, error){
	var X, X_test [][]float64
	var y, y_test []int
	if dataSetNum == 1{
		var ds utils.Dataset
		ds.LoadData()
		ds.Visualize()
		train, test := ds.TrainTestSplit(0.2)
		X, y = train.ToXY()
		X_test, y_test = test.ToXY() 
	}else if dataSetNum == 2{
		var ds utils.DatasetPasanger
		ds.LoadData()
		ds.Visualize()
		train, test := ds.TrainTestSplit(0.2)
		X, y = train.ToXY()
		X_test, y_test = test.ToXY() 
	}else if dataSetNum == 3{
		var ds utils.DatasetInsulin
		ds.LoadData()
		ds.Visualize()
		train, test := ds.TrainTestSplit(0.2)
		X, y = train.ToXY()
		X_test, y_test = test.ToXY() 
	}else{
		return nil, nil, nil, nil, fmt.Errorf("wrong data set number")
	}
	return X, y, X_test, y_test, nil
}
/*
Funkcja ShowTree wywołuje klasyfikację drzewem decyzyjnym
- pobiera dane z getDataset
- trenuje drzewo o maksymalnej głębokości 'MaxDepth'
- ocenia model na danych testowych, wylicza metryki i generuje diagram
*/
func ShowTree(dataSetNum int){
	X, y, X_test, y_test, err := getDataset(dataSetNum)
	utils.Must(err)

	tree := &DecisionTree{MaxDepth: 5}
	tree.Fit(X, y)
	tree.Analyze(X_test, y_test)
	fmt.Println("Zapisaywanie drzewa do pliku i wyświetlanie")
	utils.Must(tree.ToFlowchart(file_name))
	utils.ShowDiagram(file_name)
}
/*
Funkcja ToFLowchart generuje diagram
- tworzy plik i zapisuje w nim strukturę drzewa
- obsługuje rekurencyjne zapisywanie węzłów drzewa
*/
func (dt *DecisionTree) ToFlowchart(filename string) error {
	if dt.Tree == nil {
		return fmt.Errorf("tree is empty fit the tree first")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("graph TD\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	err = writeNode(file, dt.Tree, "root")
	if err != nil {
		return fmt.Errorf("failed to write tree structure: %v", err)
	}

	fmt.Printf("Flowchart saved to %s\n", filename)
	return nil
}
/*
Funkcja writeNode zapisuje węzeł drzewa
- sprawdza, czy węzeł jest liściem. Jeśli tak, zapisuje jego etykietę
- jeśli nie, zapisuje cehcę i próg podziału oraz rekurencyjnie przetwarza poddrzewa
*/
func writeNode(file *os.File, node *Node, nodeID string) error {
	if node.Label != nil {
		_, err := file.WriteString(fmt.Sprintf("%s[\"Class: %d\"]\n", nodeID, *node.Label))
		return err
	}

	leftID := fmt.Sprintf("%sL", nodeID)
	rightID := fmt.Sprintf("%sR", nodeID)
	_, err := file.WriteString(fmt.Sprintf("%s[\"Feature %d <= %.2f\"]\n", nodeID, node.Feature, node.Threshold))
	if err != nil {
		return err
	}
	
	_, err = file.WriteString(fmt.Sprintf("%s -->|True| %s\n", nodeID, leftID))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("%s -->|False| %s\n", nodeID, rightID))
	if err != nil {
		return err
	}

	err = writeNode(file, node.Left, leftID)
	if err != nil {
		return err
	}
	err = writeNode(file, node.Right, rightID)
	if err != nil {
		return err
	}
	return nil
}
/*
Funkcja Evaluate oblicza dokładność drzewa decyzyjnego na danych testowych
- wywołuje 'Predict'
- porównujeprzewidywane etykiety z rzeczywistymi
- liczy poprawne przewidywania i oblicza dokłądność
*/
func (dt *DecisionTree) Evaluate(testX [][]float64, testY []int) float64 {
	predictions := dt.Predict(testX)
	correct := 0
	for i := range testY {
		if testY[i] == predictions[i] {
			correct++
		}
	}
	return float64(correct) / float64(len(testY))
}

/*
Funkcja entropy oblicza entropię zbioru danych
- oblicza liczność każdej klasy w zbiorze
- dla każdej klasy oblicza prawdopodobieństwo
- stosuje wzór na entropię
- zwraca sumę jako wartość entropii
*/
func entropy(y []int) float64 {
	classCounts := make(map[int]int)
	for _, label := range y {
		classCounts[label]++
	}
	totalSamples := len(y)
	entropyValue := 0.0
	for _, count := range classCounts {
		probability := float64(count) / float64(totalSamples)
		entropyValue -= probability * math.Log2(probability)
	}
	return entropyValue
}

/*
Funkcja informationGain oblicza zysk informacji dla podziału danych
- oblicza entropię dla pełnego zbioru ('y')
- oblicza ważone entropie dla lewego i prawego podzbioru
- zwraca różnicę między entropią zbioru a ważoną sumą entropii podzbiorów
*/
func informationGain(y, yLeft, yRight []int) float64 {
	parentEntropy := entropy(y)
	leftWeight := float64(len(yLeft)) / float64(len(y))
	rightWeight := float64(len(yRight)) / float64(len(y))
	weightedEntropy := leftWeight*entropy(yLeft) + rightWeight*entropy(yRight)
	return parentEntropy - weightedEntropy
}

/*
Funkcja splitData dzieli dane na dwa podzbiory na podstawie cechy i progu
- iteruje przez próbki danych
- próbki, które spełnieją warunek, trafiają do lewego podzbioru
- pozostałe próbki trafiają do prawego podzbioru
*/
func splitData(X [][]float64, y []int, featureIndex int, threshold float64) ([][]float64, [][]float64, []int, []int) {
	var XLeft, XRight [][]float64
	var yLeft, yRight []int

	for i := range X {
		if X[i][featureIndex] <= threshold {
			XLeft = append(XLeft, X[i])
			yLeft = append(yLeft, y[i])
		} else {
			XRight = append(XRight, X[i])
			yRight = append(yRight, y[i])
		}
	}
	return XLeft, XRight, yLeft, yRight
}

/*
Funkcja buildTree buduje drzewo decyzyjne
- sprawdza warunki zatrzymania
 * wszystkie próbki należą do tej samej klasy
 * maksymalna głębokość drzewa została osiągnięta
- dla każdej cechy oblicza zysk informacji i wybiera najlepszy podział
- rekurencyjnie buduje lewe i prawe poddrzewo
*/
func (dt *DecisionTree) buildTree(X [][]float64, y []int, depth int) *Node {
	if len(unique(y)) == 1 || len(X) == 0 || (dt.MaxDepth > 0 && depth >= dt.MaxDepth) {
		label := mostCommon(y)
		return &Node{Label: &label}
	}

	bestFeature, bestThreshold, bestGain := -1, 0.0, -math.MaxFloat64
	for featureIndex := 0; featureIndex < len(X[0]); featureIndex++ {
		thresholds := uniqueFloats(column(X, featureIndex))
		for _, threshold := range thresholds {
			_, _, yLeft, yRight := splitData(X, y, featureIndex, threshold)
			if len(yLeft) > 0 && len(yRight) > 0 {
				gain := informationGain(y, yLeft, yRight)
				if gain > bestGain {
					bestFeature = featureIndex
					bestThreshold = threshold
					bestGain = gain
				}
			}
		}
	}

	if bestGain == -math.MaxFloat64 {
		label := mostCommon(y)
		return &Node{Label: &label}
	}

	XLeft, XRight, yLeft, yRight := splitData(X, y, bestFeature, bestThreshold)
	leftSubtree := dt.buildTree(XLeft, yLeft, depth+1)
	rightSubtree := dt.buildTree(XRight, yRight, depth+1)

	return &Node{
		Feature:   bestFeature,
		Threshold: bestThreshold,
		Left:      leftSubtree,
		Right:     rightSubtree,
	}
}

/*
Funkcja fit trenuje drzewo na danych treningowych
- wywołuje funkcję 'buildTree
- przechowuje wytrenowane drzewo w polu 'tree'
*/
func (dt *DecisionTree) Fit(X [][]float64, y []int) {
	dt.Tree = dt.buildTree(X, y, 0)
}

/*
Funkcja predictSample przewiduje etykietę klasy dla jednej próbki
- sprawdza, czy bieżący węzeł jest liściem. Jeśli tak to zwraca etykietę
- jeśli wartość cechy próbki jest mniejsza lub równa progowi w węźle, rekurencyjnie przechodzi do lewego poddrzewa
- w przeciwnym razie przechodzi do prawego poddrzewa
*/
func (dt *DecisionTree) predictSample(x []float64, node *Node) int {
	if node.Label != nil {
		return *node.Label
	}
	if x[node.Feature] <= node.Threshold {
		return dt.predictSample(x, node.Left)
	}
	return dt.predictSample(x, node.Right)
}

/*
Funkcja Predict przewiduje etykiety klas dla zbioru danych
- iteruje przez każdą próbkę w danych
- wywołuje predictSample dla każdej próbki
- zwraca listę przewidywanych etykiet klas
*/
func (dt *DecisionTree) Predict(X [][]float64) []int {
	predictions := make([]int, len(X))
	for i, x := range X {
		predictions[i] = dt.predictSample(x, dt.Tree)
	}
	return predictions
}

//
// Funkcje Pomocnicze 
//


/*
Funkcja unique znajduje unikalne wartości w tablicy liczb całkowitych
- wykorzystuje mapę do identyfikacji unikalnych wartości
- iteruje przez tablicę i dodaje niepowtarzające się elementy do listy
*/
func unique(arr []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/*
Funkcja uniqueFloats znajduje unikalne wartości w tablicy liczb zmiennoprzecinowych
- wykorzystuje mapę do identyfikacji unikalnych wartości
- zwraca listę unikalnych wartości jako wynik
*/
func uniqueFloats(arr []float64) []float64 {
	keys := make(map[float64]bool)
	var list []float64
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/*
Funkcja column zwraca jedną kolumnę z macierzy danych
- iteruje przez wiersze macierzy
- zbiera wartości w danej kolumnie do tablicy wynikowej
*/
func column(matrix [][]float64, index int) []float64 {
	col := make([]float64, len(matrix))
	for i := range matrix {
		col[i] = matrix[i][index]
	}
	return col
}

/*
Funkcja mostCommon znajduje najczęściej występującą wartość w tablicy 
- liczy wystąpienia każdej wartości w tablicy
- zwraca wartość o największej liczbie wystąpień
*/
func mostCommon(arr []int) int {
	counts := make(map[int]int)
	for _, value := range arr {
		counts[value]++
	}
	var maxCount, mostCommonValue int
	for key, count := range counts {
		if count > maxCount {
			maxCount = count
			mostCommonValue = key
		}
	}
	return mostCommonValue
}
/*
Funkcja Analyze oblicza i wyświetla metryki ewaluacyjne dl adrzewa decyzyjnego
- wywołuje predict aby uzyskać przewidywane etykiety klas
- oblicza metryki
- wyświetla wyniki w konsoli
*/
func (dt *DecisionTree) Analyze(X_test [][]float64, y_test []int) {
    predictions := dt.Predict(X_test)
    TP, FP, TN, FN := 0, 0, 0, 0

    for i := 0; i < len(y_test); i++ {
        if y_test[i] == 1 && predictions[i] == 1 {
            TP++
        } else if y_test[i] == 1 && predictions[i] == 0 {
            FN++
        } else if y_test[i] == 0 && predictions[i] == 1 {
            FP++
        } else if y_test[i] == 0 && predictions[i] == 0 {
            TN++
        }
    }

    accuracy := float64(TP+TN) / float64(len(y_test))
    precision := float64(TP) / float64(TP+FP)
    recall := float64(TP) / float64(TP+FN)
    f1 := 2 * (precision * recall) / (precision + recall)

    specificity := float64(TN) / float64(TN+FP)
    balancedAccuracy := (recall + specificity) / 2
    mccNumerator := float64(TP*TN - FP*FN)
    mccDenominator := math.Sqrt(float64((TP+FP)*(TP+FN)*(TN+FP)*(TN+FN)))
    mcc := mccNumerator / mccDenominator

    fmt.Printf("\nEvaluation Metrics:\n")
    fmt.Printf("--------------------\n")
    fmt.Printf("Accuracy: %.4f\n", accuracy)
    fmt.Printf("Precision: %.4f\n", precision)
    fmt.Printf("Recall (Sensitivity): %.4f\n", recall)
    fmt.Printf("F1-Score: %.4f\n", f1)
    fmt.Printf("Specificity: %.4f\n", specificity)
    fmt.Printf("Balanced Accuracy: %.4f\n", balancedAccuracy)
    fmt.Printf("Matthews Correlation Coefficient (MCC): %.4f\n", mcc)
}
