[![Build Status](https://travis-ci.org/songtianyi/docker-rec.svg?branch=master)](https://travis-ci.org/songtianyi/docker-rec)

# Docker Registry Event Collector(docker-rec)

## Get the source code
	go get github.com/songtianyi/docker-rec


## golang.org/x dep install
	mkdir -p $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git

## Note about certificates

Please note, that the DREC only works with HTTPS and so you must specify a
certificate and a key file. There are default files you can use for testing but
you should definitively create your own files.
