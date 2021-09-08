
## create network on docker

### create network
docker network create urlshortner-network
### show networks list
docker network ls
### inspect a network
docker network inspect urlshortner-network

## Prometheus
https://prometheus.io/docs/prometheus/latest/installation/
https://medium.com/aeturnuminc/configure-prometheus-and-grafana-in-dockers-ff2a2b51aa1d

### prometheus docker get command
docker pull prom/prometheus

### prometheus docker run command.

docker run -d --name prometheus --network urlshortner-network -p 9090:9090 -v ~/workspace/urlshortner/prometheus-config.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml 
- "~/workspace/urlshortner/prometheus-config.yml" is prometheus .yml config file path on host .
- -v parameter mounts "~/workspace/urlshortner/prometheus-config.yml" file on host as "/etc/prometheus/prometheus.yml" file on dokcer container.
- -d runs docker container as deamon

## grafana 
https://grafana.com/docs/grafana/latest/installation/docker/

### grafana docker run command
docker run --name grafana -d --network urlshortner-network -p 3000:3000 grafana/grafana

## mongodb

installation instruction:
https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/

### connection string
mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000

### run command
mongod --config /usr/local/etc/mongod.conf

run as background:
mongod --config /usr/local/etc/mongod.conf --fork

run as macos service:
brew services start mongodb-community
stop:
brew services stop mongodb-community

### monog caveats
==> Installing mongodb/brew/mongodb-community
==> Caveats
To have launchd start mongodb/brew/mongodb-community now and restart at login:
  brew services start mongodb/brew/mongodb-community
Or, if you don't want/need a background service you can just run:
  mongod --config /usr/local/etc/mongod.conf
==> Caveats
==> mongodb-community
To have launchd start mongodb/brew/mongodb-community now and restart at login:
  brew services start mongodb/brew/mongodb-community
Or, if you don't want/need a background service you can just run:
  mongod --config /usr/local/etc/mongod.conf

## urlshortner docker

### build
docker build -t urlshortner .

### run 
docker run --name urlshortner -d --network urlshortner-network -p 8081:8081 urlshortner

## 
### run bash shell on a running container as root
docker exec -it --user root grafana bash

### show image list
docker image ls
### remove an omage
docker image rm urlshortner

### show container list
docker container ls
### remove container
docker remove urlshortner



