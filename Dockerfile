ARG GO_VERSION=1.21.5
ARG ALPINE_VERSION=3.19
ARG XX_VERSION=1.2.1


FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
COPY --from=xx / /
RUN apk add --no-cache bash coreutils file git
ENV GO111MODULE=auto
ENV CGO_ENABLED=0
WORKDIR /src

FROM base AS version
ARG PKG=github.com/peterfromthehill/davhttpd
RUN --mount=target=. \
  VERSION=$(git describe --match 'v[0-9]*' --dirty='.m' --always --tags) REVISION=$(git rev-parse HEAD)$(if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi); \
  echo "-X ${PKG}/version.version=${VERSION#v} -X ${PKG}/version.revision=${REVISION} -X ${PKG}/version.mainpkg=${PKG}" | tee /tmp/.ldflags; \
  echo -n "${VERSION}" | tee /tmp/.version;

FROM base AS build
ARG TARGETPLATFORM
ARG LDFLAGS="-s -w"
ARG BUILDTAGS=""
RUN --mount=type=bind,target=/src \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=target=/go/pkg/mod,type=cache \
    --mount=type=bind,source=/tmp/.ldflags,target=/tmp/.ldflags,from=version \
      set -x ; xx-go build -tags "${BUILDTAGS}" -trimpath -ldflags "$(cat /tmp/.ldflags) ${LDFLAGS}" -o /usr/bin/davhttpd . \
      && xx-verify --static /usr/bin/davhttpd


FROM scratch AS binary
COPY --from=build /usr/bin/davhttpd /

FROM base AS releaser
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
WORKDIR /work
RUN --mount=from=binary,target=/build \
    --mount=type=bind,target=/src \
    --mount=type=bind,source=/tmp/.version,target=/tmp/.version,from=version \
      VERSION=$(cat /tmp/.version) \
      && mkdir -p /out \
      && cp /build/davhttpd . \
      && tar -czvf "/out/davhttpd_${VERSION#v}_${TARGETOS}_${TARGETARCH}${TARGETVARIANT}.tar.gz" * \
      && sha256sum -z "/out/davhttpd_${VERSION#v}_${TARGETOS}_${TARGETARCH}${TARGETVARIANT}.tar.gz" | awk '{ print $1 }' > "/out/davhttpd_${VERSION#v}_${TARGETOS}_${TARGETARCH}${TARGETVARIANT}.tar.gz.sha256"

FROM scratch AS artifact
COPY --from=releaser /out /


FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates
COPY --from=binary /davhttpd /davhttpd
EXPOSE 8080
ENTRYPOINT ["/davhttpd"]
