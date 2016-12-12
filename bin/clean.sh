#!/usr/bin/env bash

p=`pwd`
for file in services/api/api services/www/www services/rate/rate; do
  echo "Removing $file"
  rm $p/$file
done
cd $p
