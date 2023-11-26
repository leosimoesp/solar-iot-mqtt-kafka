#!/bin/bash

#creating 10 mqtt clients to use into simulator
for (( c=1; c < 11; c++ ))
do
    echo "creating mqtt client$c"
    docker exec mosquitto mosquitto_passwd -b /tools/mosquitto/config/password.txt "client$c" "client$c@"
done

docker-compose down 
docker-compose up -d mosquitto