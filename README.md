[![Build Status](https://travis-ci.org/kwk/docker-registry-event-collector.svg?branch=master)](https://travis-ci.org/kwk/docker-registry-event-collector) [![GoDoc](https://godoc.org/github.com/kwk/docker-registry-event-collector?status.svg)](https://godoc.org/github.com/kwk/docker-registry-event-collector) [![](https://badge.imagelayers.io/konradkleine/docker-registry-event-collector:latest.svg)](https://imagelayers.io/?images=konradkleine/docker-registry-event-collector:latest 'Get your own badge on imagelayers.io')

# Docker Registry Event Collector(docker-rec)

## Get the source code
	go get github.com/songtianyi/docker-rec


## golang.org/x dep install
	mkdir $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git

## Note about certificates

Please note, that the DREC only works with HTTPS and so you must specify a
certificate and a key file. There are default files you can use for testing but
you should definitively create your own files.
