#! /usr/bin/env bash

rm  -f  out_*.log
killall storage
nohup ./storage -file /home/biefu/rsync_insertsql_20150311/djs/upinfo.tmp -path /home/biefu/serverdata/lzservers/116.213.193.35/M9301/userdata/ &
