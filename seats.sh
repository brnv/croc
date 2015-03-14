#!/bin/bash

function is_empty {
	for sample in $(ls ./empty_seats_samples/*); do
		diff_percentage=$(compare -quiet -metric RMSE $sample $1 NULL: 2>&1)
		if [[ $diff_percentage =~ \(0(\)|\.00) ]] ; then
			echo true
		fi
	done
	echo false
}

function uniq_id {
	echo `date +%N | md5sum | cut -f1 -d' '`
}

function get_tmp_filename {
	echo /tmp/croc-$(uniq_id).png
}

function recognize {
	seat=$(get_tmp_filename)
	convert -crop 80x40+$2+$3 $table_image $seat
	result=$(is_empty $seat)
	if [[ $result == "false" ]]; then
		echo $1
	fi
	rm $seat
}

if [ -z $1 ]; then
	window_id=`xprop | grep -P "window\ id\ #\ of\ group\ leader" | grep -o "0x.*"`
	table_image=$(get_tmp_filename)
	import -window "$window_id" $table_image
else
	table_image=$1
fi

#1 hero
recognize 2 61 286
recognize 3 61 114
recognize 4 357 48
recognize 5 655 114
recognize 6 655 286
