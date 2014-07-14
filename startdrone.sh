#!/bin/bash

#usage: ./startdrone.sh github.com/xxx/xxx
git clone "git://$1" $1

./drone build $1


