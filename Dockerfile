FROM golang:alpine AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk update && apk add --no-cache git
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./out/darkroom

FROM scratch
COPY --from=builder ./out/darkroom ./out/darkroom
COPY ./application.yaml ./application.yaml
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT ["./out/darkroom"]