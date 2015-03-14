#!/bin/bash

function samples_compare {
	for sample in $(ls ./chips_samples/*); do
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

#window_id=`xprop | grep -P "window\ id" | grep -o "0x.*"`

table_image=$1

#import -window "$window_id" $table_image

chips_identifier=$(get_tmp_filename)

convert -crop 7x15+396+406 $table_image $chips_identifier

four_digit_chips_similarity=$(compare -quiet -metric RMSE $chips_identifier ./chips_samples/4_digits.png NULL: 2>&1)

three_digit_chips_similarity=$(compare -quiet -metric RMSE $chips_identifier ./chips_samples/3_digits.png NULL: 2>&1)

if [[ $four_digit_chips_similarity == *"(0)"* ]] || [[ $four_digit_chips_similarity == *"(0.0"* ]] ; then
	echo 4 digit
elif [[ $three_digit_chips_similarity  == *"(0)"* ]] || [[ $three_digit_chips_similarity  == *"(0.0"* ]] ; then
	echo 3 digit
else
	echo 2 digit
fi
