package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/patrikeh/go-deep/training"
	"math"
	"os"
	"strconv"
)

type AlzheimerData struct {
	Age                float64
	Gender             float64
	Ethnicity          float64
	EducationLevel     float64
	BMI                float64
	Smoking            float64
	AlcoholConsumption float64
	PhysicalActivity   float64
	DietQuality        float64
	MemoryComplaints   float64
	BehavioralProblems float64
	ADL                float64
	Confusion          float64
	Disorientation     float64
	PersonalityChanges float64
	DifficultyTasks    float64
	Forgetfulness      float64
	Diagnosis          float64
}

type AlzheimerDataset []AlzheimerData

// Print wypisuje dane pacjentów. Można opcjonalnie podać liczbę rekordów do wypisania.
func (ds AlzheimerDataset) Print(optionalArgs ...int) {
	var limit int = -1
	if len(optionalArgs) == 1 {
		limit = optionalArgs[0]
	}
	for i, data := range ds {
		fmt.Printf("%+v\n", data)
		if limit > 0 && i+1 >= limit {
			break
		}
	}
}

// LoadData ładuje dane z pliku CSV do AlzheimerDataset.
// filePath to ścieżka do pliku CSV.
func (ds *AlzheimerDataset) LoadData(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("błąd otwierania pliku: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("błąd odczytu pliku CSV: %v", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		var data AlzheimerData
		data.Age, _ = strconv.ParseFloat(row[1], 64)
		data.Gender, _ = strconv.ParseFloat(row[2], 64)
		data.Ethnicity, _ = strconv.ParseFloat(row[3], 64)
		data.EducationLevel, _ = strconv.ParseFloat(row[4], 64)
		data.BMI, _ = strconv.ParseFloat(row[5], 64)
		data.Smoking, _ = strconv.ParseFloat(row[6], 64)
		data.AlcoholConsumption, _ = strconv.ParseFloat(row[7], 64)
		data.PhysicalActivity, _ = strconv.ParseFloat(row[8], 64)
		data.DietQuality, _ = strconv.ParseFloat(row[9], 64)
		data.MemoryComplaints, _ = strconv.ParseFloat(row[10], 64)
		data.BehavioralProblems, _ = strconv.ParseFloat(row[11], 64)
		data.ADL, _ = strconv.ParseFloat(row[12], 64)
		data.Confusion, _ = strconv.ParseFloat(row[13], 64)
		data.Disorientation, _ = strconv.ParseFloat(row[14], 64)
		data.PersonalityChanges, _ = strconv.ParseFloat(row[15], 64)
		data.DifficultyTasks, _ = strconv.ParseFloat(row[16], 64)
		data.Forgetfulness, _ = strconv.ParseFloat(row[17], 64)
		data.Diagnosis, _ = strconv.ParseFloat(row[18], 64)

		*ds = append(*ds, data)
	}
	return nil
}

// ToXY zwraca dane w formie macierzy cech (X) i wektora wyników (Y).
func (ds AlzheimerDataset) ToXY() ([][]float64, []float64) {
	X := make([][]float64, len(ds))
	Y := make([]float64, len(ds))
	for i, data := range ds {
		X[i] = []float64{
			data.Age, data.Gender, data.Ethnicity, data.EducationLevel,
			data.BMI, data.Smoking, data.AlcoholConsumption, data.PhysicalActivity,
			data.DietQuality, data.MemoryComplaints, data.BehavioralProblems,
			data.ADL, data.Confusion, data.Disorientation, data.PersonalityChanges,
			data.DifficultyTasks, data.Forgetfulness,
		}
		Y[i] = data.Diagnosis
	}
	return X, Y
}

// ToExamples konwertuje dane do formatu wymaganych przez bibliotekę "go-deep".
func (ds AlzheimerDataset) ToExamples() []training.Example {
	var examples []training.Example
	for _, data := range ds {
		examples = append(examples, training.Example{
			Input: []float64{
				data.Age, data.Gender, data.Ethnicity, data.EducationLevel,
				data.BMI, data.Smoking, data.AlcoholConsumption, data.PhysicalActivity,
				data.DietQuality, data.MemoryComplaints, data.BehavioralProblems,
				data.ADL, data.Confusion, data.Disorientation, data.PersonalityChanges,
				data.DifficultyTasks, data.Forgetfulness,
			},
			Response: []float64{data.Diagnosis},
		})
	}
	return examples
}

// Normalize normalizuje dane w zbiorze, przekształcając wartości cech do zakresu [0, 1]
func (ds *AlzheimerDataset) Normalize() error {
	minVals := make([]float64, 17)
	maxVals := make([]float64, 17)

	for i := range minVals {
		minVals[i] = math.MaxFloat64
		maxVals[i] = -math.MaxFloat64
	}

	for _, data := range *ds {
		features := []float64{
			data.Age, data.Gender, data.Ethnicity, data.EducationLevel,
			data.BMI, data.Smoking, data.AlcoholConsumption, data.PhysicalActivity,
			data.DietQuality, data.MemoryComplaints, data.BehavioralProblems,
			data.ADL, data.Confusion, data.Disorientation, data.PersonalityChanges,
			data.DifficultyTasks, data.Forgetfulness,
		}
		for i, val := range features {
			if val < minVals[i] {
				minVals[i] = val
			}
			if val > maxVals[i] {
				maxVals[i] = val
			}
		}
	}

	for i := range *ds {
		features := []float64{
			(*ds)[i].Age, (*ds)[i].Gender, (*ds)[i].Ethnicity, (*ds)[i].EducationLevel,
			(*ds)[i].BMI, (*ds)[i].Smoking, (*ds)[i].AlcoholConsumption, (*ds)[i].PhysicalActivity,
			(*ds)[i].DietQuality, (*ds)[i].MemoryComplaints, (*ds)[i].BehavioralProblems,
			(*ds)[i].ADL, (*ds)[i].Confusion, (*ds)[i].Disorientation, (*ds)[i].PersonalityChanges,
			(*ds)[i].DifficultyTasks, (*ds)[i].Forgetfulness,
		}
		for j, val := range features {
			if maxVals[j] != minVals[j] {
				features[j] = (val - minVals[j]) / (maxVals[j] - minVals[j])
			} else {
				features[j] = 0
			}
		}
		(*ds)[i].Age = features[0]
		(*ds)[i].BMI = features[4]
		(*ds)[i].Diagnosis = features[len(features)-1]
	}
	return nil
}
