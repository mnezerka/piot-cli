#!/usr/bin/env bash

# generate the plot
docker run --rm -v $(pwd):/work remuslazar/gnuplot -e \
 "set xlabel 'Datum'; set ylabel 'Teplota';
  set xdata time;
  set timefmt \"%Y-%m-%dT%H:%M:%S\";
  set grid;
  set key autotitle columnhead;
  set datafile separator ',';
  set term png size 800,380;
  set output 'chart.png';
  plot 'data.csv' using 1:2 with lines"

 #plot 'data.csv' every 2::0 skip 1 using 2:3 title 'Global Temperature' with lines linewidth 2;"
