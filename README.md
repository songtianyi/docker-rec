[![Build Status](https://travis-ci.org/songtianyi/docker-rec.svg?branch=master)](https://travis-ci.org/songtianyi/docker-rec)

# Docker Registry Event Collector(docker-rec)

## Get the source code
	go get github.com/songtianyi/docker-rec


## golang.org/x dep install
	mkdir -p $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git

## docker registery notification setup
    notifications:
      endpoints:
        - name: docker-rec
          disabled: false
          url: https://x.x.x.x:8080/events
          timeout: 500ms
          threshold: 5
          backoff: 1s
