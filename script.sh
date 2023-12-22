CGO_ENABLED=0 GOOS=linux go build -o ./mdns-broadcast

scp ./mdns-broadcast serverB:/root/app
