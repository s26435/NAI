package models

import (
	"venv/utils"

	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func ShowInsuliNN(){
	var dataset utils.DatasetInsulin
	dataset.LoadData()
	dataset.MustNormalize()
	train, test := dataset.TrainTestSplit(0.2)
	train_data := train.ToExamples()
	test_data := test.ToExamples()
	n := deep.NewNeural(&deep.Config{
		Inputs: 3,
		Layout: []int{3, 16, 32, 64, 32, 16, 1},
		Activation: deep.ActivationSigmoid,
		Mode: deep.ModeBinary,
		Weight: deep.NewNormal(1.0, 0.0),
		Bias: true,

	})

	optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
	trainer := training.NewTrainer(optimizer, 50)
	trainer.Train(n, train_data, test_data, 100) 
	utils.EvaluateModelInsulin(test_data, n)
}


