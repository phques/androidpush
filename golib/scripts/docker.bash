#!/bin/sh

#exec sudo docker run -ti -v $GOPATH/src:/src -v ~/code/mppq:/mppq golang/mobile /bin/bash -i
exec sudo docker run -ti -v $GOPATH/src:/src golang/mobile /bin/bash -i

