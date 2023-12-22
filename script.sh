CGO_ENABLED=0 GOOS=linux go build -o ./mdns-broadcast

scp serverB:/root/mdns ./mdns-broadcast 
