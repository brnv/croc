#!/bin/bash

# refactor
# implement in go

window_id=`xprop | grep -P "window\ id" | grep -o "0x.*"`

uniq_name=/tmp/screenshot_`date +%N | md5sum | cut -f1 -d' '`

import -window "$window_id" $uniq_name.png

left_card=$uniq_name.left.png

right_card=$uniq_name.right.png

convert -crop 46x30+346+340 $uniq_name.png $left_card

convert -crop 46x30+396+340 $uniq_name.png $right_card

left_is_uniq="1"

right_is_uniq="1"

for f in $(ls ./card_samples/*/*); do
	diff_percentage=$(compare -quiet -metric RMSE $left_card $f NULL: 2>&1)
	if [[ $diff_percentage == *"(0)"* ]]; then
		left_is_uniq="0"
		echo $f | rev | cut -f1-2 -d'/' | rev | tr -d / | cut -f1 -d'.'
	fi

	diff_percentage=$(compare -quiet -metric RMSE $right_card $f NULL: 2>&1)
	if [[ $diff_percentage == *"(0)"* ]]; then
		right_is_uniq="0"
		echo $f | rev | cut -f1-2 -d'/' | rev | tr -d / | cut -f1 -d'.'
	fi
done

if [[ $left_is_uniq == "1" ]]; then
    read -p "enter left card value: " value
    mkdir -p ./samples/$value
    read -p "enter left card suit: " suit
    mv $left_card ./samples/$value/$suit.png
fi

if [[ $right_is_uniq == "1" ]]; then
    read -p "enter right card value: " value
    mkdir -p ./samples/$value
    read -p "enter right card suit: " suit
    mv $right_card ./samples/$value/$suit.png
fi
