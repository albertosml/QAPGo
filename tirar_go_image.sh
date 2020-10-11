#!/bin/bash

OPTION=$1
PWD=$2
DEST=$3

# Build or run image
if [ "$OPTION" == "build" ]; then
	# Build image
	podman build -t my_go_image .
elif [ "$OPTION" == "run" ]; then
	# Take default local directory if we don't have it
        if [ -z "$PWD" ]; then
            PWD=~/Desktop/qap/go
        fi

	# Take default dest if we don't have it
	if [ -z "$DEST" ]; then
        	DEST=/home/ubuntu
	fi

	# Run go container
	podman run --rm -it --name my_go_container -v $PWD:$DEST my_go_image
else
	echo "You must specify 'build' or 'run' option"
fi

# Cheatsheet
# Build -> go build file -o program_name | ./program_name
# Run -> go run file
