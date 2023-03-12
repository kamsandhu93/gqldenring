FROM golang:1.18-buster as builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server


FROM debian:buster-slim
RUN set -x && apt-get update

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server
#COPY config.json /

# Run the web service on container startup.
CMD ["/app/server"]