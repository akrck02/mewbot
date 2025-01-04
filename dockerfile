FROM golang:1.23-bookworm AS base

# LABEL about the custom image
LABEL maintainer="akrck02@gmail.com"
LABEL version="0.2"
LABEL description="This is a custom Docker Image for mewbot execution"

RUN mkdir -p /home/app/mewbot
RUN touch /home/app/mewbot/.env

ENV DISCORD_BOT_TOKEN="none"

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build -o mewbot


# Start the application
CMD ["/build/mewbot"]


