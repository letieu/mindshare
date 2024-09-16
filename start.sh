#!/bin/bash

PID=$(ps -ef | grep mindshare | grep -v grep | awk '{print $2}' | head -n 1)
if [ -n "$PID" ]; then
  kill "$PID"
else
  echo "No mindshare process found."
fi

rm -rf mindshare static templates
tar -xvf deploy.tar.gz
nohup ./mindshare &> mindshare.log &
