#!/bin/sh
docker run --rm -it -v `pwd`:/go/src/github.com/bitly/oauth2_proxy -w /go/src/github.com/bitly/oauth2_proxy haoling/go-buildtools /go/src/github.com/bitly/oauth2_proxy/dist.sh
