#!/bin/bash
#i=50
for ((i=72;i>0;i--))
do
da=$(date "+%Y-%m-%d" --date="-$i day")
echo $da >> date.log
done
