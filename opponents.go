package main

import "strconv"

//@TODO: move to config
var (
	betSamples          = "/tmp/croc/bet_digits/*"
	betCompareThreshold = 0.2
	betDigitWidth       = 8
	betDigitHeight      = 12
)

type Opponent struct {
	Index      int
	Raiser     bool
	Limper     bool
	ChipsInPot int
}

// visgrep offsets
var opponentsOffets = map[int]string{
	1: "501,36 -1",
	2: "657,92 -1",
	3: "700,205 -1",
	4: "609,321 -1",
	5: "hero position, not used",
	6: "157,321 -1",
	7: "65,205 -1",
	8: "108,92 -1",
	9: "264,36 -1",
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

func (opponent Opponent) GetImageSnippets() []ImageSnippet {
	return []ImageSnippet{
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Index][0],
			OffsetY: limpDigitOffsetsY[opponent.Index],
		},
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Index][1],
			OffsetY: limpDigitOffsetsY[opponent.Index],
		},
		ImageSnippet{
			Width:   betDigitWidth,
			Height:  betDigitHeight,
			OffsetX: limpDigitOffsetsX[opponent.Index][2],
			OffsetY: limpDigitOffsetsY[opponent.Index],
		},
	}
}

func (table *Table) OpponentsRecognize() {
	offsets := getSubimageManyOffsets(
		table.Image.Path,
		"/tmp/croc/opponents_hand",
	)

	for _, offset := range offsets {
		for opponentIndex, opponentOffset := range opponentsOffets {
			if opponentOffset == offset {
				table.Opponents = append(table.Opponents, Opponent{
					Index: opponentIndex,
				})
			}
		}
	}

	for index, opponent := range table.Opponents {
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

			table.Opponents[index].ChipsInPot += betInteger
			table.Opponents[index].Limper = true
		}

	}
}

func (table *Table) RaisersRecognize() {
	for index, opponent := range table.Opponents {
		if opponent.ChipsInPot == 0 {
			table.Opponents[index].Raiser = true
		}
	}
}
