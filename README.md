
# Usage

## Test TCP Listener 
``` sh
go run ./cmd/tcplistener | tee /tmp/rawget.http
```

### GET Request
``` sh
curl http://localhost:42069/coffee 
```

### POST Request
``` sh
curl -X POST -H "Content-Type: application/json" -d '{"flavor":"dark mode"}' http://localhost:42069/coffee
```

## Test UDP 
``` sh
go run ./cmd/udpsender
nc -u -l 42069
```