#!/bin/sh

#alias shosetRun='go run -race test/*.go 5'

sleep 2

./bin/shoset_build 5 B 0 0 rien &
#P2=$!

wait