package main

import (
	"strconv"
	"strings"
)

//@TODO: move to config
var (
	potDigitSamples     = "/tmp/croc/pot_digits/*"
	potTypeSamples      = "/tmp/croc/pot_types/*"
	potTypeWidth        = 14
	potTypeHeight       = 13
	potTypeOffsetX      = 385
	potOffsetY          = 154
	potDigitWidth       = 9
	potDigitHeight      = 13
	potCompareThreshold = 0.05
)

type Pot struct {
	Number
}

func (table *Table) PotRecognize() {
	pot := Pot{
		Number: Number{
			Digits: []ImageSnippet{},
		},
	}

	potType, err := table.Image.NumberTypeRecognize(
		ImageSnippet{
			Width:   potTypeWidth,
			Height:  potTypeHeight,
			OffsetX: potTypeOffsetX,
			OffsetY: potOffsetY,
		},
		potTypeSamples,
		potCompareThreshold,
	)

	if err != nil {
		table.Errors = append(table.Errors, err.Error())
		return
	}

	switch potType {
	case "3":
		pot.Number.Digits = pot.GetPotImageSnippets(
			[]int{401, 413, 422},
		)
	}

	potSize := ""

	for _, potDigit := range pot.Number.Digits {
		digit, err := recognize(
			table.Image.Crop(potDigit),
			potDigitSamples,
			potCompareThreshold,
		)

		if err != nil {
			table.Errors = append(table.Errors, err.Error())
			return
		}
		potSize += digit
	}

	table.Pot, _ = strconv.Atoi(strings.TrimLeft(potSize, "0"))
}

func (pot Pot) GetPotImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		potDigitWidth,
		potDigitHeight,
		potOffsetY,
		offsets,
	)
}
