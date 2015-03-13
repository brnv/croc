#!/bin/bash

window_id=`xprop | grep -P "window\ id" | grep -o "0x.*"`

uniq_name=/tmp/screenshot_`date | md5sum | cut -f1 -d' '`

import -window "$window_id" $uniq_name.png

pot_identifier=$uniq_name.pot.id.png

convert -crop 8x13+390+154 $uniq_name.png $pot_identifier

two_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/2_digits.png NULL: 2>&1)
three_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/3_digits.png NULL: 2>&1)

if [[ $two_digit_pot_similarity == *"(0)"* ]]; then
	echo its two-digit pot
	exit
fi

if [[ $three_digit_pot_similarity == *"(0)"* ]]; then
	echo its three-digit pot
	exit
fi

echo its four-digit pot
