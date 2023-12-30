############################
# STEP 1 build an executable
############################
# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

#setup env
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
# RUN apk add --no-cache git bash

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
# WORKDIR /app/blog
# COPY . .

WORKDIR /src/htmx-blog
COPY . .
# COPY go.mod go.sum /src/netbanking
# Install library dependencies
# RUN go mod download
RUN go mod verify

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
# COPY . /go/src/project/
RUN GOOS=linux GOARCH=amd64 GIN_MODE=release go build -o app

############################
# STEP 2 build a small image
############################
# This results in a single layer image
FROM alpine:latest
# LABEL MAINTAINER Amrendra <amren1254@gmail.com>
RUN apk --no-cache add ca-certificates 

# Add new user 'appuser'
# RUN adduser -D appuser
# USER appuser

COPY --from=builder /src/htmx-blog/app /blogapp

# WORKDIR /home/appuser/app/blog

# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /src/htmx-blog/.env.local /go/bin/app
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=root
ENV DB_HOSTNAME=localhost
ENV DB_PORT=5432
ENV DB_NAME=postgres
ENV DB_SSLMODE=disable
ENV HOST=localhost
ENV PORT=8081
ENV STATIC_FILE_DIR=./static
# COPY --from=build /src/htmx-blog/app /go/bin/app
# RUN ls
# USER 405

EXPOSE 8081

# CMD [ "./blogapp" ]
ENTRYPOINT ["/blogapp"]
