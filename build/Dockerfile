FROM golang:1.20-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
RUN apk update && apk add build-base

COPY . .

RUN GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-w -s" -o darkroom main.go

FROM alpine

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/darkroom ./darkroom

RUN chmod +x ./darkroom
ENV PORT 3000
EXPOSE 3000

CMD ["./darkroom", "server"]
