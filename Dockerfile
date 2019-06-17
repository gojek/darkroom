FROM golang:alpine AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk update && apk add --no-cache git openssh-client
RUN git config --global url."git@***REMOVED***:".insteadOf "https://***REMOVED***/"
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 400 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan ***REMOVED*** >> /root/.ssh/known_hosts
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o darkroom cmd/darkroom/main.go
RUN rm /root/.ssh/id_rsa

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /app/darkroom ./darkroom
RUN chmod +x ./darkroom
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT ["./darkroom"]
