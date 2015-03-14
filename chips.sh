#!/bin/bash

function samples_compare {
	for sample in $(ls ./pot_samples/{0..9}{,_l}.png 2>/dev/null); do
		diff_percentage=$(compare -quiet -metric RMSE $sample $1 NULL: 2>&1)
		if [[ $diff_percentage =~ \(0(\)|\.0) ]] ; then
			basename $sample | cut -f1 -d'.' | cut -f1 -d'_' | cat | tr -d '\n'
			return
		fi
	done
	echo new
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
	result=$(samples_compare $digit)
	if [[ $result == "new" ]]; then
		read -p "enter digit value: " value
		mv $digit ./pot_samples/$value.png
		return
	fi
	rm $digit
	echo -n $result
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

if [[ $four_digit_chips_similarity =~ \(0(\)|\.0|\.1) ]] ; then
	recognize 401
	recognize 413
	recognize 422
	recognize 431
elif [[ $three_digit_chips_similarity =~ \(0(\)|\.0|\.1) ]] ; then
	recognize 407
	recognize 416
	recognize 425
else
	recognize 412
	recognize 421
fi

