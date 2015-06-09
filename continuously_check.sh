#!/bin/sh

make sources | entr -cr sh -c 'make || exit; pkill -9 jwpack; sleep .1; ./jwpack -l 0.0.0.0:5000'
