# Project of home dashboard #

## Author: Vikerian - Vilem Kebrt ##

## Version: 0.02 - no release yet ##

### License

### "THE BEER-WARE LICENSE" (Revision 42)

### <kebrtv@gmail.com> wrote this file

### As long as you retain this notice you can do whatever you want with this stuff

### If we meet some day, and you think this stuff is worth it, you can buy me a beer in return

### Vikerian

#### Description

- uses some web clients to get data
- reads data from home mqtt queues
- store data into mongodb (web data),
  siridb (time slice data),
  valkey.io (instead of recently closed redis,
  session data and such data which are not persistent).

VALKEY.IO

- Fork and continuer of redis:
- docker image: valkey/valkey:8.1.1-alpine3.21

SiriDB:

- Timeline database written in C++
  - Run (build) from : ./ContainerFiles/SiriDB-Containerfile

MongoDB:

- NOSQL document storage database

---

Run mongo on podman:

podman pull docker.io/mongodb/mongodb-community-server:latest
podman run --detach --name mongo -p 27017:27017 -v ./mongo-data:/data/db docker.io/mongodb/mongodb-community-server:latest

Run valkey on podman:
podman run -d -p "6379:6379" --name valkey "valkey/valkey:8.1.1-alpine3.21"
