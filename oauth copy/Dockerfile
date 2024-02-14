# Multi-stage build to compile the Go application


FROM golang:latest AS buildoauth 

WORKDIR /application

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-oauth

# Build a minimal image for running the application
FROM alpine:3.11.3

WORKDIR /

COPY --from=buildoauth /docker-oauth .

EXPOSE 8080

ENTRYPOINT ["/docker-oauth"]


