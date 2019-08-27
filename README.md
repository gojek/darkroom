# Darkroom - Yet Another Image Proxy

<p align="center"><img src="website/static/img/darkroom-logo.png" width="360"></p>

[![build status](https://travis-ci.com/gojek/darkroom.svg?branch=master)](https://travis-ci.com/gojek/darkroom)
[![Coverage Status](https://coveralls.io/repos/github/gojek/darkroom/badge.svg?branch=master)](https://coveralls.io/github/gojek/darkroom?branch=master)
[![Docs latest](https://img.shields.io/badge/Docs-latest-blue.svg)](https://www.gojek.io/darkroom/)
[![GoDoc](https://godoc.org/github.com/gojek/darkroom?status.svg)](https://godoc.org/github.com/gojek/darkroom)
[![Go Report Card](https://goreportcard.com/badge/github.com/gojek/darkroom)](https://goreportcard.com/report/github.com/gojek/darkroom)
[![GolangCI](https://golangci.com/badges/github.com/gojek/darkroom.svg)](https://golangci.com)
[![GitHub Release](https://img.shields.io/github/release/gojek/darkroom.svg?style=flat)](https://github.com/gojek/darkroom/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

## Introduction

[Darkroom](https://www.gojek.io/darkroom/) combines the [storage backend](pkg/storage) and the [image processor](pkg/processor) and acts as an `Image Proxy` on your image source.  
You may [implement](https://www.gojek.io/darkroom/docs/customization#custom-storage-example) your own `Storage` and `Processor` interfaces to gain custom functionality while still keeping other Darkroom Server functionality.  
The native implementations focus on speed and resiliency.

## Features

Darkroom supports several image operations which are documented [here](https://www.gojek.io/darkroom/docs/usage/size).

## Installation

```bash
go get -u github.com/gojek/darkroom
```
Other ways to run can be found [here](https://www.gojek.io/darkroom/docs/getting-started#running-the-image-proxy-service).

### Contributing Guide

Read our [contributing guide](./CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to Darkroom.

## License

Darkroom is [MIT licensed](./LICENSE).
