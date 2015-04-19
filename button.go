package main

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
			buttonWidth, buttonHeight, 510, 117,
		},

		2: ImageSnippet{
			buttonWidth, buttonHeight, 605, 164,
		},

		3: ImageSnippet{
			buttonWidth, buttonHeight, 632, 200,
		},

		4: ImageSnippet{
			buttonWidth, buttonHeight, 539, 311,
		},

		5: ImageSnippet{
			buttonWidth, buttonHeight, 458, 329,
		},

		6: ImageSnippet{
			buttonWidth, buttonHeight, 252, 318,
		},

		7: ImageSnippet{
			buttonWidth, buttonHeight, 146, 204,
		},

		8: ImageSnippet{
			buttonWidth, buttonHeight, 176, 164,
		},

		9: ImageSnippet{
			buttonWidth, buttonHeight, 281, 117,
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

		table.ButtonPosition = seat

		break
	}
}
