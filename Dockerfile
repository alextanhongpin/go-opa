FROM golang:latest

COPY server.go server.go

CMD ["go", "run", "server.go"]
