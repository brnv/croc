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
	1: 115,
	2: 145,
	3: 250,
	4: 289,
	5: 301,
	6: 289,
	7: 247,
	8: 144,
	9: 126,
}

var limpDigitOffsetsX = map[int]map[int]int{
	1: map[int]int{0: 445, 1: 457, 2: 466},
	2: map[int]int{0: 535, 1: 547, 2: 556},
	3: map[int]int{0: 577, 1: 589, 2: 598},
	4: map[int]int{0: 514, 1: 526, 2: 535},
	5: map[int]int{0: 393, 1: 405, 2: 414},
	6: map[int]int{0: 261, 1: 273, 2: 282},
	7: map[int]int{0: 392, 1: 204, 2: 213},
	8: map[int]int{0: 235, 1: 247, 2: 256},
	9: map[int]int{0: 337, 1: 349, 2: 358},
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

func (table *Table) LimpersRecognize() {
	for index := 1; index <= 9; index++ {
		table.Limpers = append(table.Limpers, Limper{
			Position: index,
		})
	}

	for index, opponent := range table.Limpers {
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

			betInteger, _ := strconv.Atoi(recognized)

			table.Limpers[index].BetSize += betInteger
		}

	}
}
