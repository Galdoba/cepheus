#!/bin/bash

cd ../cmd
for d in */; do
    (cd "$d" && echo go install)
done
cd ../scripts
