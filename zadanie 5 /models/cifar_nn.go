package models

import (
	
	"venv/utils"

	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func ShowCifarNN() {
	var dataset utils.CifarDataSet
	dataset.LoadData()
	dataset.Normalize()

	train, test := dataset.TrainTestSplit(0.2)
	train_data := train.ToExamples()
	test_data := test.ToExamples()

	n := deep.NewNeural(&deep.Config{
		Inputs:     3072,
		Layout:     []int{64, 32, 16, 10},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(0.01, 0.0),
		Bias:       true,
	})

	optimizer := training.NewAdam(0.001, 0.9, 0.999, 1e-8)
	trainer := training.NewBatchTrainer(optimizer, 1, 256, 4) // optimizer, verbose, batchsize, num workwrs
	trainer.Train(n, train_data, test_data, 1)
	utils.EvaluateModelCifar(test_data, n, 10, utils.CifarLabels)
}
