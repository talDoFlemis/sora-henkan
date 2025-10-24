#!/bin/bash

set -e

cd /opt/sora-henkan

/usr/local/bin/docker-compose up -d --pull always
echo "Sora Henkan application started"