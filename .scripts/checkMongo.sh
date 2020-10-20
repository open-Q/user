#!/usr/bin/env bash

while ! curl -s localhost:27017/ > curlAnswer.txt;
     do sleep 1
done

rm curlAnswer.txt