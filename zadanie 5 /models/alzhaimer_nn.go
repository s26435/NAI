package models

import (
	"fmt"
	"math"
	"math/rand"
	"venv/utils"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func ShowAlzhaimer() {
	
	var dataset utils.AlzheimerDataset
	err := dataset.LoadData("alzheimers_disease_data.csv")
	utils.Must(err)
	for i, record := range dataset {
		if math.IsNaN(record.Age) || math.IsNaN(record.BMI) || math.IsNaN(record.Diagnosis) {
			fmt.Printf("Błędne dane w wierszu %d: %+v\n", i, record)
			return
		}
	}

	err = dataset.Normalize()
	utils.Must(err)

	data := dataset.ToExamples()
	trainingData, testData := splitData(data, 0.8)

	config := deep.Config{
		Inputs:     len(trainingData[0].Input),
		Layout:     []int{64, 32, 16, 1},
		Activation: deep.ActivationTanh,
		Mode:       deep.ModeRegression,
		Weight:     deep.NewNormal(0.6, 0.1),
		Bias:       true,
	}
	network := deep.NewNeural(&config)

	trainer := training.NewTrainer(training.NewSGD(0.0001, 0.1, 1e-6, true), 50)
	trainer.Train(network, trainingData, testData, 2000)

	utils.EvaluateAlzhaimerModel(network, testData)
}

func splitData(data []training.Example, trainRatio float64) (training.Examples, training.Examples) {
	rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })

	splitIndex := int(float64(len(data)) * trainRatio)
	return data[:splitIndex], data[splitIndex:]
}
