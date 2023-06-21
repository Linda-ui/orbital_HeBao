#!/bin/bash


# This script starts the hertz gateway and the nacos server. 
# It should be executed before the backend servers startup script.


# start the nacos server (the nacos directory should be located at the home direcotory)
bash ~/nacos/bin/startup.sh -m standalone

# start the hertz gateway (at background)
go run ./hertz_gateway 

# shut down the nacos server when the gateway shuts down
# bash ~/nacos/bin/shutdown.sh