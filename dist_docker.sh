#!/bin/sh
docker run --rm -it -v `pwd`:/app -w /app haoling/go-buildtools /app/dist.sh
