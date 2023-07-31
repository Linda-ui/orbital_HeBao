#!/bin/bash

# DO NOT EDIT.
# This script starts all the backend kitex servers.

# function to recursively find .go files recursively and execute them
function backend_server_startup() {
    directory="$1"
    for subdir in "$directory"/*; do
        if [[ -d "$subdir" ]]; then
            if [[ -f "$subdir/main.go" ]]; then
                # start the backend kitex server (at background)
                go run "$subdir/main.go" &
            fi
        fi
    done
}

# starting directory of the backend servers
server_directory="./kitex_services"

# start the backend kitex servers
backend_server_startup $server_directory

# wait for all the background processes to finish
wait
