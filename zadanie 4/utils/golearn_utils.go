package utils

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"strings"
	"strconv"
)

func makeAtrr() *[]base.Attribute {
	atrTable := []string{
		"FixedAcidity",
		"VolatileAcidity",
		"CitricAcid",
		"ResidualSugar",
		"Chlorides",
		"FreeSulfurDioxide",
		"TotalSulfurDioxide",
		"Density",
		"PH",
		"Sulphates",
		"Alcohol",
		"Quality",
	}
	attrs := make([]base.Attribute, 12)
	for i, name := range atrTable {
		attrs[i] = new(base.FloatAttribute)
		attrs[i].SetName(name)
		attrs[i].GetSysValFromString("1.0")
	}

	return &attrs
}


func makeData(ds Dataset) [][]float64 {
	newData := make([][]float64, 0, len(ds))
	for _, data := range ds {
		temp := []float64{
			data.FixedAcidity,
			data.VolatileAcidity,
			data.CitricAcid,
			data.ResidualSugar,
			data.Chlorides,
			data.FreeSulfurDioxide,
			data.TotalSulfurDioxide,
			data.Density,
			data.PH,
			data.Sulphates,
			data.Alcohol,
			data.Quality,
		}
		newData = append(newData, temp)
	}
	return newData
}

func MakeInastances(ds Dataset) *base.DenseInstances {
	attrs := *makeAtrr()
	instances := makeData(ds)

	newInst := base.NewDenseInstances()
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = newInst.AddAttribute(a)
	}
	fmt.Println(newSpecs)
	newInst.Extend(len(instances))

	for j := 0; j < len(instances)-1; j++ {
		for i := 0; i < len(attrs); i++ {
			newInst.Set(newSpecs[i], j, newSpecs[i].GetAttribute().GetSysValFromString(strings.TrimSpace(strconv.FormatFloat(instances[j][i], 'f', -1, 64))))
		}
	}

	for i := 0; i < len(attrs); i++ {
		if attr, ok := newInst.AllAttributes()[i].(*base.FloatAttribute); !ok {
			panic("Invalid cast")
		} else {
			attr.Precision = 4
		}
	}
	newInst.AddClassAttribute(attrs[len(attrs)-1])
	newInst.Set(newSpecs[11], 0, newSpecs[11].GetAttribute().GetSysValFromString(strings.TrimSpace(strconv.FormatFloat(instances[0][11], 'f', -1, 64))))
	return newInst
}
