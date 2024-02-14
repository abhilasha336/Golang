
# builder image
FROM golang:alpine as builder

RUN mkdir /build
WORKDIR /build

COPY go.mod /build/
COPY go.sum /build/
COPY . /build/

RUN go mod download

RUN go build -o base.cmd /build/main.go


# generate clean, final image for end users
FROM alpine:3.11.3

COPY --from=builder /build/base.cmd .

EXPOSE 8080

# executable
ENTRYPOINT [ "./base.cmd"]
