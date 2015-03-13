#!/bin/bash

# sorry for that code

window_id=`xprop | grep -P "window\ id" | grep -o "0x.*"`

uniq_name=/tmp/screenshot_`date +%N | md5sum | cut -f1 -d' '`

import -window "$window_id" $uniq_name.png

pot_identifier=$uniq_name.pot.id.png

convert -crop 8x13+390+154 $uniq_name.png $pot_identifier

two_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/2_digits.png NULL: 2>&1)
three_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/3_digits.png NULL: 2>&1)

if [[ $two_digit_pot_similarity == *"(0)"* ]]; then
	first_is_new="1"
	second_is_new="1"

	first_digit=$uniq_name.first.png
	second_digit=$uniq_name.second.png

	convert -crop 9x13+403+154 $uniq_name.png $first_digit
	convert -crop 9x13+412+154 $uniq_name.png $second_digit

	for f in $(ls ./pot_samples/*); do
		diff_percentage=$(compare -quiet -metric RMSE $first_digit $f NULL: 2>&1)
		if [[ $diff_percentage == *"(0)"* ]] || [[ $diff_percentage == *"(0.00"* ]] ; then
			first_is_new="0"
			basename $f | cut -f1 -d'.' | cat | tr -d '\n'
		fi
	done

	for f in $(ls ./pot_samples/*); do
		diff_percentage=$(compare -quiet -metric RMSE $second_digit $f NULL: 2>&1)
		if [[ $diff_percentage == *"(0)"* ]] || [[ $diff_percentage == *"(0.00"* ]] ; then
			second_is_new="0"
			basename $f | cut -f1 -d'.' | cat | tr -d '\n'
			echo
		fi
	done


	if [[ $first_is_new == "1" ]]; then
		read -p "enter first digit: " digit
		mv $first_digit ./pot_samples/$digit.png
	fi

	if [[ $second_is_new == "1" ]]; then
		read -p "enter second digit: " digit
		mv $second_digit ./pot_samples/$digit.png
	fi

	exit
fi

if [[ $three_digit_pot_similarity == *"(0)"* ]]; then
	echo its three-digit pot
	exit
fi

echo its four-digit pot
