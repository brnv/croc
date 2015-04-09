package main

import "strconv"

//@TODO: move to config
var (
	betSamples          = "/tmp/croc/bet_digits/*"
	betCompareThreshold = 0.2
	betDigitWidth       = 8
	betDigitHeight      = 12
)

type Limper struct {
	Position int
	BetSize  int
}

var limpDigitOffsetsY = map[int]int{
	1: 301,
	2: 289,
	3: 247,
	4: 144,
	5: 126,
	6: 115,
	7: 145,
	8: 250,
	9: 289,
}

var limpDigitOffsetsX = map[int]map[int]int{
	1: map[int]int{0: 393, 1: 405, 2: 414},
	2: map[int]int{0: 261, 1: 273, 2: 282},
	3: map[int]int{0: 392, 1: 204, 2: 213},
	4: map[int]int{0: 235, 1: 247, 2: 256},
	5: map[int]int{0: 337, 1: 349, 2: 358},
	6: map[int]int{0: 445, 1: 457, 2: 466},
	7: map[int]int{0: 535, 1: 547, 2: 556},
	8: map[int]int{0: 577, 1: 589, 2: 598},
	9: map[int]int{0: 514, 1: 526, 2: 535},
}

func (opponent Limper) GetImageSnippets() []ImageSnippet {
	return []ImageSnippet{
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Position][0],
			OffsetY: limpDigitOffsetsY[opponent.Position],
		},
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Position][1],
			OffsetY: limpDigitOffsetsY[opponent.Position],
		},
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Position][2],
			OffsetY: limpDigitOffsetsY[opponent.Position],
		},
	}
}

func (table Table) LimpersRecognize() []Limper {
	opponents := []Limper{}
	for index := 1; index <= 9; index++ {
		opponents = append(opponents, Limper{
			Position: index,
		})
	}

	for index, opponent := range opponents {
		for _, betDigit := range opponent.GetImageSnippets() {
			img := table.Image.Crop(betDigit)

			recognized, err := recognize(
				img,
				betSamples,
				betCompareThreshold,
			)

			if err != nil {
				continue
			}

			recognizedInt, _ := strconv.Atoi(recognized)

			opponents[index].BetSize += recognizedInt
		}

	}

	return opponents
}
