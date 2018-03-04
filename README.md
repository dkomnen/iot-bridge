# iot-bridge

## Setup

`docker run -it -p 1883:1883 -p 9001:9001 eclipse-mosquitto`

From project root:

`go run cmd/bifrost/main.go --connnect-address="http://CONNECT_HOST:PORT"

## Running Valkyrie

After everything is up and running:

From project root:

`go run cmd/valkyrie/main.go temp --tick INTERVAL -l LOWER_TEMP_LIMIT -u UPPER_TEMP_LIMIT`





