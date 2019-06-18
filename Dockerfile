FROM golang:alpine AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk update && apk add --no-cache git
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o darkroom cmd/darkroom/main.go

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /app/darkroom ./darkroom
RUN chmod +x ./darkroom
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT ["./darkroom"]
