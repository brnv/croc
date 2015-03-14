#!/bin/bash

# implement in go

function samples_compare {
	for sample in $(ls ./card_samples/*/*); do
		diff_percentage=$(compare -quiet -metric RMSE $sample $1 NULL: 2>&1)
		if [[ $diff_percentage == *"(0)"* ]] || [[ $diff_percentage == *"(0.00"* ]] ; then
			echo $sample | rev | cut -f1-2 -d/ | cut -f2 -d. | tr -d '/' | rev | tr -d '\n'
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
	card=$(get_tmp_filename)
	convert -crop 46x30+$1+340 $table_image $card
	samples_compare $card
	rm $card
}

if [ -z $1 ]; then
	window_id=`xprop | grep -P "window\ id\ #\ of\ group\ leader" | grep -o "0x.*"`
	table_image=$(get_tmp_filename)
	import -window "$window_id" $table_image
else
	table_image=$1
fi

recognize 346
recognize 396
echo
