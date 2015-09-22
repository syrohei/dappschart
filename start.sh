#!/bin/sh

nohup ./appusd &
nohup ./mailserver &
nohup /home/ubuntu/.go/bin/btcd &
