package models

import (
	"fmt"
	"math"
	"zad4/utils"
)
/*
Funkcja ShowSVM to główna funkcja wywołująca klasyfikacje SVM:
- pobiera dane z getDataset
- uczy model SVM na danych treningowych
- analizuje wyniki klasyfikacji na danych testowych
*/
func ShowSVM(dataSetNum int){
	X, y, X_test, y_test, err := getDataset(dataSetNum)
	utils.Must(err)
	uniqueClasses := unique(y)

	mcSVM := NewMultiClassSVM(0.001, 0.01, 1000, uniqueClasses)
	mcSVM.Fit(X, y)

	mcSVM.Analyze(X_test, y_test)
}
// Struktura dla klasyfikatora SVM
type SVM struct {
	LearningRate float64
	LambdaParam  float64
	NIters       int
	W            []float64
	B            float64
}
// Funkcja tworzy nowy obiekt klasyfikatora SVM
func NewSVM(learningRate, lambdaParam float64, nIters int) *SVM {
	return &SVM{
		LearningRate: learningRate,
		LambdaParam:  lambdaParam,
		NIters:       nIters,
		W:            nil,
		B:            0,
	}
}
/* 
Funkcja fit trenuke model SVM
- inicjalizuje wektory wag 'W' i przesunięcie 'B'
- iteracyjnie optymalizuje 'W' i 'B' na podstawie warunku marginu (Hinge Loss):
  * jeśli punkt jest poprawnie sklasyfikowany, zmniejsza regularizację
  * w przeciwnym wypadku aktualizuj wagi i przesunięcie
- powtarza proces przez zadaną liczbę iteracji ('NIters')
*/
func (svm *SVM) fit(X [][]float64, y []int) {
	nSamples := len(X)
	nFeatures := len(X[0])

	svm.W = make([]float64, nFeatures)

	y_ := make([]int, nSamples)
	for i := range y {
		if y[i] <= 0 {
			y_[i] = -1
		} else {
			y_[i] = 1
		}
	}

	for i := 0; i < svm.NIters; i++ {
		for idx, x_i := range X {
			var dotProduct float64
			for j := 0; j < nFeatures; j++ {
				dotProduct += x_i[j] * svm.W[j]
			}
			condition := float64(y_[idx])*(dotProduct+svm.B) >= 1
			if condition {
				for j := 0; j < nFeatures; j++ {
					svm.W[j] -= svm.LearningRate * (2 * svm.LambdaParam * svm.W[j])
				}
			} else {
				for j := 0; j < nFeatures; j++ {
					svm.W[j] -= svm.LearningRate * (2*svm.LambdaParam*svm.W[j] - float64(y_[idx])*x_i[j])
				}
				svm.B -= svm.LearningRate * float64(y_[idx])
			}
		}
	}
}
// Funkcja predict przewiduje wyniki dla zbioru danych
func (svm *SVM) predict(X [][]float64) []float64 {
	nSamples := len(X)
	predictions := make([]float64, nSamples)

	for i, x := range X {
		var approx float64
		for j, feature := range x {
			approx += feature * svm.W[j]
		}
		approx += svm.B

		predictions[i] = approx
	}

	return predictions
}
// Struktura dla wieloklasowego SVM
type MultiClassSVM struct {
	SVMs       []*SVM
	UniqueClasses []int
}
/*
Funkcja NewMultiClass tworzy nowy model wieloklasowego SVM
*/
func NewMultiClassSVM(learningRate, lambdaParam float64, nIters int, uniqueClasses []int) *MultiClassSVM {
	svms := make([]*SVM, len(uniqueClasses))
	for i := range uniqueClasses {
		svms[i] = NewSVM(learningRate, lambdaParam, nIters)
	}
	return &MultiClassSVM{
		SVMs:       svms,
		UniqueClasses: uniqueClasses,
	}
}
/*
Klasa Fit uczy model wieloklasowego SVM
- dla każdej klasy tworzy model binarny
- uczy każdy model na odpowiednich
*/
func (mc *MultiClassSVM) Fit(X [][]float64, y []int) {
	for i, class := range mc.UniqueClasses {
		yBinary := make([]int, len(y))
		for j, label := range y {
			if label == class {
				yBinary[j] = 1
			} else {
				yBinary[j] = -1
			}
		}

		mc.SVMs[i].fit(X, yBinary)
	}
}
/*
Funkcja Predict przewiduje etykiety klas dla danych wejściowych
- iteruje przez wszystkie próbki w danych wejściowych
- dla każdej próbki oblicza wynik dla każdego modelu SVM
- porównuje wyniki i wybiera klasę z najwyższym wynikiem
- zwraca listę przewidywanych etykiet klas
*/
func (mc *MultiClassSVM) Predict(X [][]float64) []int {
	nSamples := len(X)
	predictions := make([]int, nSamples)

	for i, x := range X {
		maxScore := math.Inf(-1)
		bestClass := -1

		for j, svm := range mc.SVMs {
			score := 0.0
			for k, feature := range x {
				score += feature * svm.W[k]
			}
			score += svm.B

			if score > maxScore {
				maxScore = score
				bestClass = mc.UniqueClasses[j]
			}
		}

		predictions[i] = bestClass
	}

	return predictions
}
/*
Funkcja Evaluate oblicza dokładność modelu na danych testowych
- wywołuje fukncję Predict
- porównuje przewidywane etykiety z rzeczywistymi etykietami
- liczy liczbę poprawnych przewidywań
- oblicza i zwraca dokładność jako procent poprawnych przewidywań
*/
func (mc *MultiClassSVM) Evaluate(X [][]float64, y []int) float64{
	predictions := mc.Predict(X)
	nSamples := len(y)
	correct := 0

	for i := 0; i < nSamples; i++ {
		if predictions[i] == y[i] {
			correct++
		}
	}

	return float64(correct) / float64(nSamples) * 100
}
/*
Funkcja Analyze oblicza metryki jakości klasyfikacji dla każdego modelu
- oblicza metryki takie jak dokładność, precyzja, czułość itp
- wyświetla wyniki w konsoli
*/
func (mc *MultiClassSVM) Analyze(X [][]float64, y []int) {
    predictions := mc.Predict(X)
    uniqueLabels := unique(y)

    // Metryki dla każdej klasy
    metrics := make(map[int]map[string]float64)
    for _, class := range uniqueLabels {
        TP, FP, TN, FN := 0, 0, 0, 0

        for i := 0; i < len(y); i++ {
            if y[i] == class && predictions[i] == class {
                TP++
            } else if y[i] == class && predictions[i] != class {
                FN++
            } else if y[i] != class && predictions[i] == class {
                FP++
            } else {
                TN++
            }
        }

        accuracy := float64(TP+TN) / float64(TP+TN+FP+FN)
        precision := float64(TP) / float64(TP+FP)
        recall := float64(TP) / float64(TP+FN)
        f1 := 2 * (precision * recall) / (precision + recall)
        specificity := float64(TN) / float64(TN+FP)

        metrics[class] = map[string]float64{
            "Accuracy":    accuracy,
            "Precision":   precision,
            "Recall":      recall,
            "F1-Score":    f1,
            "Specificity": specificity,
        }
    }

    // Ogólne metryki
    overallAccuracy := mc.Evaluate(X, y) / 100
    fmt.Printf("\nOverall Accuracy: %.4f\n", overallAccuracy)
    fmt.Println("Metrics per class:")

    for class, metric := range metrics {
        fmt.Printf("\nClass %d:\n", class)
        fmt.Printf("  Accuracy: %.4f\n", metric["Accuracy"])
        fmt.Printf("  Precision: %.4f\n", metric["Precision"])
        fmt.Printf("  Recall: %.4f\n", metric["Recall"])
        fmt.Printf("  F1-Score: %.4f\n", metric["F1-Score"])
        fmt.Printf("  Specificity: %.4f\n", metric["Specificity"])
    }
}
