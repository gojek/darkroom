---
id: getting-started
title: Getting Started
---



## Installation
```bash
go get -u github.com/gojek/darkroom
```



## Running the Image Proxy Service
The project has docker images available. They can be tested locally or can be be deployed to production.

Create a file containing the environment variables mentioned in [`config.example.yaml`](./config.example.yaml) and save it as `config.env`
> Note: Bucket credentials are dummy, you need to provide your own credentials.
```bash
DEBUG=true
LOG_LEVEL=debug

APP_NAME=darkroom
APP_VERSION=0.0.1
APP_DESCRIPTION="A realtime image processing service"

SOURCE_KIND=s3
SOURCE_BUCKET_NAME=bucket-name
SOURCE_BUCKET_REGION=bucket-region
SOURCE_BUCKET_ACCESSKEY=AKIA*************
SOURCE_BUCKET_SECRETKEY=4y/*******************************
SOURCE_PATHPREFIX=/uploads

PORT=3000

CACHE_TIME=31536000

SOURCE_HYSTRIX_COMMANDNAME=S3_ADAPTER
SOURCE_HYSTRIX_TIMEOUT=5000
SOURCE_HYSTRIX_MAXCONCURRENTREQUESTS=100
SOURCE_HYSTRIX_REQUESTVOLUMETHRESHOLD=10
SOURCE_HYSTRIX_SLEEPWINDOW=10
SOURCE_HYSTRIX_ERRORPERCENTTHRESHOLD=25
```
Build the docker image and run it with the config created.
```bash
make docker-image
docker run -p 80:3000 --env-file ./config.env ${USER}/darkroom:latest
```
