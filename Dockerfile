FROM golang:1.26.6-bookworm

ARG MT_MULTISERVER_PROXY_REPO=fondazione-golinelli/mt-multiserver-proxy
ARG VERSION
ARG BUILDPLATFORM
ARG BUILDARCH
ARG TARGETARCH

ENV GONOSUMCHECK=github.com/HimbeerserverDE/mt-multiserver-proxy
ENV GONOSUMDB=github.com/HimbeerserverDE/mt-multiserver-proxy
ENV GOPRIVATE=github.com/HimbeerserverDE/mt-multiserver-proxy
ENV GOFLAGS=-trimpath

COPY . /go/src/github.com/HimbeerserverDE/mt-multiserver-proxy

RUN mkdir /usr/local/mt-multiserver-proxy
RUN git config --system url."https://github.com/${MT_MULTISERVER_PROXY_REPO}".insteadOf "https://github.com/HimbeerserverDE/mt-multiserver-proxy"
# Build the checked-out tree, including an upstream merge not yet pushed to GitHub.
# Native target-platform builds also keep CGO enabled for SQLite and Go plugins.
WORKDIR /go/src/github.com/HimbeerserverDE/mt-multiserver-proxy
RUN PROXY_BUILD_VERSION="${VERSION:-$(TZ=UTC git show -s --abbrev=12 --date=format-local:%Y%m%d%H%M%S --format=v0.0.0-%cd-%h)}" && \
    GOBIN=/usr/local/mt-multiserver-proxy go install \
    -ldflags "-X github.com/HimbeerserverDE/mt-multiserver-proxy.buildVersion=${PROXY_BUILD_VERSION}" ./cmd/...
WORKDIR /usr/local/mt-multiserver-proxy

VOLUME ["/usr/local/mt-multiserver-proxy"]

EXPOSE 40000/udp

CMD ["/usr/local/mt-multiserver-proxy/mt-multiserver-proxy"]
