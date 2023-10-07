FROM golang:1.21

ENV PROJECT_DIR=/usr/src/events \
    CGO_ENABLED=0

WORKDIR /usr/src/events
RUN mkdir "/build"
COPY . .
RUN go install github.com/githubnemo/CompileDaemon@latest
ENTRYPOINT CompileDaemon -build="go build -o ./build/app ./cmd/main.go" -command="./build/app api"
