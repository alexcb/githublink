VERSION 0.7
FROM alpine:3.16

deps:
    FROM golang:1.19-alpine3.16
    RUN apk add --update --no-cache \
        bash \
        bash-completion \
        binutils \
        ca-certificates \
        coreutils \
        curl \
        findutils \
        g++ \
        git \
        grep \
        less \
        make \
        openssl \
        util-linux

    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.0
    WORKDIR /code

    # otherwise, this would be needed
    COPY go.mod go.sum .
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY --dir cmd ./
    SAVE IMAGE

lint:
    FROM +code
    COPY ./.golangci.yaml ./
    RUN golangci-lint run

githublink:
    FROM +code
    ARG RELEASE_TAG="dev"
    ARG GOOS
    ARG GO_EXTRA_LDFLAGS
    ARG GOARCH
    RUN test -n "$GOOS" && test -n "$GOARCH"
    ARG GOCACHE=/go-cache
    RUN mkdir -p build
    ENV CGO_ENABLED=0
    RUN --mount=type=cache,target=$GOCACHE \
        go build \
            -o build/githublink \
            -ldflags "-X main.Version=$RELEASE_TAG $GO_EXTRA_LDFLAGS" \
            cmd/main.go
    SAVE ARTIFACT build/githublink AS LOCAL "build/$GOOS/$GOARCH/githublink"

githublink-darwin-amd64:
    COPY \
        --build-arg GOOS=darwin \
        --build-arg GOARCH=amd64 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +githublink/githublink /build/githublink
    SAVE ARTIFACT /build/githublink AS LOCAL "build/darwin/amd64/githublink"

githublink-darwin-arm64:
    COPY \
        --build-arg GOOS=darwin \
        --build-arg GOARCH=arm64 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +githublink/githublink /build/githublink
    SAVE ARTIFACT /build/githublink AS LOCAL "build/darwin/arm64/githublink"

githublink-linux-amd64:
    COPY \
        --build-arg GOOS=linux \
        --build-arg GOARCH=amd64 \
        --build-arg GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static" \
        +githublink/githublink /build/githublink
    SAVE ARTIFACT /build/githublink AS LOCAL "build/linux/amd64/githublink"

githublink-linux-arm64:
    COPY \
        --build-arg GOOS=linux \
        --build-arg GOARCH=arm64 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +githublink/githublink /build/githublink
    SAVE ARTIFACT /build/githublink AS LOCAL "build/linux/arm64/githublink"

githublink-all:
    PIPELINE
    TRIGGER push pipetest
    TRIGGER pr pipetest
    BUILD +githublink-linux-amd64
    BUILD +githublink-linux-arm64
    BUILD +githublink-darwin-amd64
    BUILD +githublink-darwin-arm64


release:
    FROM node:13.10.1-alpine3.11
    RUN npm install -g github-release-cli@v1.3.1
    WORKDIR /release
    COPY +githublink-linux-amd64/githublink ./githublink-linux-amd64
    COPY +githublink-linux-arm64/githublink ./githublink-linux-arm64
    COPY +githublink-darwin-amd64/githublink ./githublink-darwin-amd64
    COPY +githublink-darwin-arm64/githublink ./githublink-darwin-arm64
    ARG --required RELEASE_TAG
    ARG EARTHLY_GIT_HASH
    ARG BODY="No details provided"
    RUN --secret GITHUB_TOKEN=+secrets/GITHUB_TOKEN test -n "$GITHUB_TOKEN"
    RUN --push \
        --secret GITHUB_TOKEN=+secrets/GITHUB_TOKEN \
        github-release upload \
        --owner alexcb \
        --repo githublink \
        --commitish "$EARTHLY_GIT_HASH" \
        --tag "$RELEASE_TAG" \
        --name "$RELEASE_TAG" \
        --body "$BODY" \
        ./githublink-*
