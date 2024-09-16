#!/bin/bash

go build
tar -czf deploy.tar.gz mindshare static templates
scp deploy.tar.gz tieu@$HOST:/home/tieu/mindshare

rm deploy.tar.gz
rm mindshare
