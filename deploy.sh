#!/bin/bash
docker run --rm --privileged \
  -v $(pwd):/darkroom \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(which docker):$(which docker) \
  -w /darkroom \
  -e GITHUB_TOKEN \
  -e DOCKER_USER \
  -e DOCKER_PASSWORD \
  bepsays/ci-goreleaser \
  docker login docker.io -u $DOCKER_USER -p $DOCKER_PASSWORD && \
  goreleaser release --skip-validate --rm-dist