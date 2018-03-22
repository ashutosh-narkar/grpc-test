# To build for linux:
# env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o grpc-proxy

# To build for Mac:
# env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o grpc-proxy-darwin


FROM ubuntu:16.04
WORKDIR /app
COPY grpc-proxy /app/
COPY config.json /app/
EXPOSE 50051

# Install curl
RUN apt-get update && apt-get install -y curl

# Install vim
RUN ["apt-get", "install", "-y", "vim"]

ENTRYPOINT ["./grpc-proxy"]
