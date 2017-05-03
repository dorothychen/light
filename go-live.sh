#!/usr/bin/env bash

# ticker key is just 209 on localhost
url="feelingcolor.herokuapp.com/ctrl/$1"
if [$ticker_key = ""]; then
    url="localhost:9090/ctrl/$1/"
    ticker_key="209"
fi

echo $url$ticker_key
curl $url$ticker_key
