


## Prometheus
https://prometheus.io/docs/prometheus/latest/installation/
https://medium.com/aeturnuminc/configure-prometheus-and-grafana-in-dockers-ff2a2b51aa1d

### prometheus docker get command
docker pull prom/prometheus

### prometheus docker run command.

docker run -d --name prometheus -p 9090:9090 -v ~/workspace/urlshortner/prometheus-config.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml 
- "~/workspace/urlshortner/prometheus-config.yml" is prometheus .yml config file path on host .
- -v parameter mounts "~/workspace/urlshortner/prometheus-config.yml" file on host as "/etc/prometheus/prometheus.yml" file on dokcer container.
- -d runs docker container as deamon

## grafana 
https://grafana.com/docs/grafana/latest/installation/docker/

### grafana docker run command
docker run --name grafana -d -p 3000:3000 grafana/grafana

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

###
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

