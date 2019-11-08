# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/Generalbelly/q-torial-api
COPY . .

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o q-torial-api

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/Generalbelly/q-torial-api/q-torial-api /q-torial-api

# Run the web service on container startup.
CMD ["/q-torial-api"]
