#! /usr/bin/bash

loc=/mnt/dfathom/data/2100/FLOOD_MAP-1_3ARCSEC-NW_OFFSET-1in1000-PLUVIAL-DEFENDED-DEPTH-2100-SSP5_8.5-PERCENTILE50-v3.1/
outfile=max_depths.txt
i=0
max=0

# echo 'max_depths' > $outfile

for file in $loc*
do
  currmax=$(gdalinfo -mm $file | grep 'Computed' | awk 'BEGIN{FS=","} {print $2}' | awk 'BEGIN{FS="."} {print $1}' )
  if [ $currmax -gt $max ]; then
    max=$currmax
  fi
done

echo Max is $max