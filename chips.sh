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
	convert -crop 9x13+$1+407 $table_image -colorspace Gray $digit
	samples_compare $digit
	rm $digit
}

if [ -z $1 ]; then
	window_id=`xprop | grep -P "window\ id\ #\ of\ group\ leader" | grep -o "0x.*"`
	table_image=$(get_tmp_filename)
	import -window "$window_id" $table_image
else
	table_image=$1
fi

chips_identifier=$(get_tmp_filename)

convert -crop 7x15+396+406 $table_image $chips_identifier

four_digit_chips_similarity=$(compare -quiet -metric RMSE $chips_identifier ./pot_samples/4_digits_chips.png NULL: 2>&1)

three_digit_chips_similarity=$(compare -quiet -metric RMSE $chips_identifier ./pot_samples/3_digits_chips.png NULL: 2>&1)

if [[ $four_digit_chips_similarity == *"(0)"* ]] || [[ $four_digit_chips_similarity == *"(0.0"* ]] ; then
	recognize 401
	recognize 413
	recognize 422
	recognize 431
elif [[ $three_digit_chips_similarity  == *"(0)"* ]] || [[ $three_digit_chips_similarity  == *"(0.0"* ]] ; then
	recognize 407
	recognize 416
	recognize 425
else
	recognize 412
	recognize 421
fi

