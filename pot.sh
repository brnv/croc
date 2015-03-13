#!/bin/bash

function samples_compare {
	for sample in $(ls ./pot_samples/*); do
		diff_percentage=$(compare -quiet -metric RMSE $sample $1 NULL: 2>&1)
		if [[ $diff_percentage == *"(0)"* ]] || [[ $diff_percentage == *"(0.0"* ]] ; then
			basename $sample | cut -f1 -d'.' | cat | tr -d '\n'
		fi
	done
}

function uniq_id {
	echo `date +%N | md5sum | cut -f1 -d' '`
}

function get_tmp_filename {
	echo /tmp/croc-$(uniq_id).png
}

function recognize {
	digit=$(get_tmp_filename)
	convert -crop 9x13+$1+154 $table_image $digit
	samples_compare $digit
	rm $digit
}

window_id=`xprop | grep -P "window\ id" | grep -o "0x.*"`

table_image=$(get_tmp_filename)

import -window "$window_id" $table_image

pot_identifier=$(get_tmp_filename)

convert -crop 14x13+390+154 $table_image $pot_identifier

two_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/2_digits.png NULL: 2>&1)

three_digit_pot_similarity=$(compare -quiet -metric RMSE $pot_identifier ./pot_samples/3_digits.png NULL: 2>&1)

if [[ $two_digit_pot_similarity == *"(0)"* ]] || [[ $two_digit_pot_similarity == *"(0.0"* ]] ; then
	recognize 407
	recognize 416
elif [[ $three_digit_pot_similarity  == *"(0)"* ]] || [[ $three_digit_pot_similarity  == *"(0.0"* ]] ; then
	recognize 403
	recognize 412
	recognize 421
else
	recognize 397
	recognize 409
	recognize 418
	recognize 427
fi
