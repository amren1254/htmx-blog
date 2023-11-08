############################
# STEP 1 build an executable
############################
# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build

#setup env
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git bash

# Create appuser
# ENV USER=appuser
# ENV UID=10001

# RUN addgroup --system appuser

# RUN adduser \  
#     --system\
#     --disabled-password \    
#     --gecos "" \    
#     --home "/nonexistent" \    
#     --shell "/sbin/nologin" \    
#     --no-create-home \    
#     --uid "10001" \    
#     "appuser"

# List project dependencies with go.mod and go.sum
# These layers are only re-built when Gopkg files are updated
WORKDIR /src/htmx-blog
COPY . .
# COPY go.mod go.sum /src/netbanking
# Install library dependencies
RUN go mod download
RUN go mod verify

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
# COPY . /go/src/project/
RUN GOOS=linux GOARCH=amd64 go build -o app

############################
# STEP 2 build a small image
############################
# This results in a single layer image
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /src/htmx-blog/app /go/bin/app
USER 405

EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]
