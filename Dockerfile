# Linux container coming with go ROOT configured at /go.
FROM golang:1.15

# Add the source code to the container.
# We copy all source code because the client depends on some miner packages.
ADD . /go/src/

# CD into the source code directory.
WORKDIR /go/src/bazo-client

# Build the application.
RUN go build -o /bazo-client

# Define the start command when this container is run.
CMD ["/bazo-client", "rest"]