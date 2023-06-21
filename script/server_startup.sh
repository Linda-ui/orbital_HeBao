#!/bin/bash


# This script starts all the backend kitex servers.
# It should be executed after the gateway startup script.

# function to recursively find .go files recursively and execute them
function backend_server_startup() {
    directory="$1"
    for file in $directory/*; do
        if [[ -d "$file" ]]; then
            backend_server_startup "$file"
        elif [[ -f "$file" ]]; then
            if [[ "${file##*.}" == "go" ]]; then
                # start the backend kitex server (at background)
                go run "$file" &
            fi
        fi
    done
}

# starting directory of the backend servers
server_directory="./kitex_services/kitex_server"

# start the backend kitex servers
backend_server_startup $server_directory

# wait for all the background processes to finish
wait
