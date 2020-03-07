# Stage 1 : Build
# Use golang base image
FROM golang:stretch as builder

# Declare args
ARG REVISION
ARG RELEASE_TAG

# Fresh ssl certs
RUN apt-get update && apt-get install -y ca-certificates

# Create src dir
RUN mkdir /opt/gambling-bot

# Set working directory
WORKDIR /opt/gambling-bot

# Deps
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Change workdir
WORKDIR /opt/gambling-bot/cmd

# build binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o gambling-bot

# Stage 2 : run !
FROM debian:stretch

RUN useradd -ms /bin/bash gamble

# Get certs
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# image-spec annotations using labels
# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.source="https://github.com/Namarand/gambling-bot"
LABEL org.opencontainers.image.revision=${GIT_COMMIT_SHA}
LABEL org.opencontainers.image.version=${RELEASE_TAG}
LABEL org.opencontainers.image.authors="Wilfried OLLIVIER"
LABEL org.opencontainers.image.title="gambling-bot"
LABEL org.opencontainers.image.description="gambling-bot runtime"

# Copy our static executable.
COPY --from=builder /opt/gambling-bot/cmd/gambling-bot /opt/gambling-bot

USER gamble

ENV LOGLEVEL=WARN

# Run the binary
ENTRYPOINT ["/opt/gambling-bot"]
