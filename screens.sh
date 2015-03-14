#!/bin/bash

window_id=`xprop | grep -P "window\ id\ #\ of" | grep -o "0x.*"`

while true; do
	uniq_name=/tmp/screenshot_`date +%N | md5sum | cut -f1 -d' '`
	import -window "$window_id" $uniq_name.png
	sleep 1s
done
