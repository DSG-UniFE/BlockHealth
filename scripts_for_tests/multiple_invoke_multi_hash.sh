#!/bin/bash

# Qui ci interessa più valutare la dimensione di una transazione che cambia all'aumentare di hash

echo "start script"
date


MS=0

#max 8 hash (+2 default)
while [ $MS -lt 9 ]; do
	#echo 'name: ' $(($MS % 2))
	node invoke.js "$MS" "$MS" > ./log/log-invoke-dynamic-hash-${MS}.log 2>&1 & 
	echo "ms= $MS"
	MS=$((MS+1))
	sleep 3
done

echo "finish script"
date

