package main

import (
	"strconv"
)

//@TODO: move to config
var (
	buttonSamples          = "/tmp/croc/button/*"
	buttonCompareThreshold = 0.05
	buttonWidth            = 16
	buttonHeight           = 13
)

var button map[int]ImageSnippet

func (table *Table) ButtonRecognize() {
	button := map[int]ImageSnippet{
		1: ImageSnippet{
			buttonWidth, buttonHeight, 458, 329,
		},

		2: ImageSnippet{
			buttonWidth, buttonHeight, 252, 318,
		},

		3: ImageSnippet{
			buttonWidth, buttonHeight, 146, 204,
		},

		4: ImageSnippet{
			buttonWidth, buttonHeight, 176, 164,
		},

		5: ImageSnippet{
			buttonWidth, buttonHeight, 281, 117,
		},

		6: ImageSnippet{
			buttonWidth, buttonHeight, 510, 117,
		},

		7: ImageSnippet{
			buttonWidth, buttonHeight, 605, 164,
		},

		8: ImageSnippet{
			buttonWidth, buttonHeight, 632, 200,
		},

		9: ImageSnippet{
			buttonWidth, buttonHeight, 539, 311,
		},
	}

	for seat, buttonPosition := range button {
		_, err := recognize(
			table.Image.Crop(buttonPosition),
			buttonSamples,
			buttonCompareThreshold,
		)

		if err != nil {
			continue
		}

		table.Button = strconv.Itoa(seat)
	}
}
