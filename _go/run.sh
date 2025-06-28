#!/bin/sh
#
export MQTT_URL="mqtt://test.mosquitto.org:1883"
export MQTT_AUTH="dXNlcjpwYXNzCg=="
export VALKEY_URL="tcp://127.0.0.1:6379"
export VALKEY_AUTH="dXNlcjpwYXNzCg=="
export SIRIDB_URL="dGNwOi8vMTI3LjAuMC4xOjkwMDAK"
export SIRIDB_AUTH="tcp://127.0.0.1:9000"
export MONGO_URL="mongodb://127.0.0.1:27017"
#export MONGO_AUTH="dGNwOi8vMTI3LjAuMC4xOjkwMDA"

/usr/bin/go run src/main.go
