#!/bin/bash

# Qui ci interessa più valutare le prestazioni con un intenso traffico di dati (caso pessimo = 8 hash). 1TPS per esempio

echo "start script"
date


MS=0

#max 8 hash (+2 default)
while [ $MS -lt 10 ]; do
	#echo 'name: ' $(($MS % 2))
	node invoke.js "$MS" "8" > ./log/log-invoke-1s-${MS}.log 2>&1 & 
	echo "ms= $MS"
	MS=$((MS+1))
	sleep 1
done

echo "finish script"
date

