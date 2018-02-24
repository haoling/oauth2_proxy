#!/bin/sh
docker run --rm -it -v `pwd`:/go/src/github.com/bitly/oauth2_proxy -w /go/src/github.com/bitly/oauth2_proxy golang bash -c 'curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && /go/src/github.com/bitly/oauth2_proxy/dist.sh'
