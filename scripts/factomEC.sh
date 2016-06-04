#!/bin/bash

chain="$(echo "one two three four, this is a test of making an entry.  More to test, Test, Test, test" $RANDOM $RANDOM $RANDOM | factom-cli mkchain -e $RANDOM -e $RANDOM e1 | gawk "{print \$3}")"
sleep 201
echo "Made Chain: " $chain

for i in `seq 1 10`;
do
    entry="$(echo "Here are some random numbers: " $RANDOM $RANDOM $RANDOM $RANDOM $RANDOM $RANDOM)"
    echo $entry | factom-cli put -c $chain e1
    echo $entry
done