#!/bin/bash

window_id=`xwininfo | grep -Po "Window\sid:\s([\S.]+)\s" | cut -f3 -d' '`

if [[ -z $window_id ]]; then
	exit
fi

hand_num=1

while true; do
	echo -n \#

	start=`date +%s%N`
	result=`./croc --wid $window_id -v $1`
	code=$?
	stop=`date +%s%N`

	if [[ $code == 0 ]]; then
		echo $hand_num $(($stop - $start)) $result

		hand_num=$(($hand_num+1))
	fi

	echo $result | grep -Poq MANUAL
	if [[ $? == 0 ]]; then
		 espeak "Attention!"
	fi

	sleep 2s
done
