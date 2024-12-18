package models

import (
	"venv/utils"
	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)


//funkcja pokazująca i porównująca działanie dwóch różnych sieci neuronowych
// z takimi samymi parametrami, ale o różnej wielkości 
func ShowFashionNN(){
	var fds utils.FashionDataset
	fds.LoadData()
	fds.Normalize()
	
	test, train := fds.TrainTestSplit(0.2)
	train_data := train.ToExamples()
	test_data := test.ToExamples()

	n1 := deep.NewNeural(&deep.Config{
		Inputs:     784,
		Layout:     []int{128, 64, 32, 10}, 
		Activation: deep.ActivationReLU,
		Mode:       deep.ModeMultiClass, 
		Weight:     deep.NewNormal(0.5, 0.0),
		Bias:       true,
	})

	n2 := deep.NewNeural(&deep.Config{
		Inputs:     784,
		Layout:     []int{64, 32, 16, 10}, 
		Activation: deep.ActivationReLU,
		Mode:       deep.ModeMultiClass, 
		Weight:     deep.NewNormal(0.5, 0.0),
		Bias:       true,
	})

	labels := utils.Labels

	optimizer := training.NewAdam(0.001, 0.9, 0.999, 1e-8)

	trainer1 := training.NewBatchTrainer(optimizer, 1, 256, 10) // optimizer, verbose, batchsize, num workwrs
	trainer2 := training.NewBatchTrainer(optimizer, 1, 256, 10) // optimizer, verbose, batchsize, num workwrs

	trainer1.Train(n1, train_data, test_data, 10)
	trainer2.Train(n2, train_data, test_data, 10)

	utils.EvaluateModelCifar(test_data, n1, 10, labels)
	utils.EvaluateModelCifar(test_data, n2, 10, labels)
}
