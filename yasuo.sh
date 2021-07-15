#!/bin/bash
for line in `cat date.log`
do
    #tar -xzvf "$line.tar.gz"
    #tar -zcvf "$line.tar.gz" $line
    #mv "$line.tar.gz" /home/yangkan/rxy/data/
    sh copy_data.sh $line
done
