# Build stage
FROM golang:1.23.4-bookworm AS builder

WORKDIR /src
COPY . . 

ARG CGO_ENABLED=0
ENV CGO_ENABLED ${CGO_ENABLED}

ARG VERSION="snapshot"
ENV VERSION ${VERSION}

ARG BUILD_DEBUG
ENV BUILD_DEBUG ${BUILD_DEBUG}

RUN mkdir -p /build/usr/bin                                         \
 && ./build.sh -b                                                   \
 && mv bin/unwise /build/usr/bin/                                   \
 && ldd /build/usr/bin/unwise | tr -s '[:blank:]' '\n'              \
  | grep '^/'                                                       \
  | xargs -I % sh -c 'mkdir -p /build/$(dirname %); cp % /build/%;' 

# Final stage
FROM scratch

COPY --from=builder --chown=1000:1000 /build /

USER 1000:1000

EXPOSE 3123

CMD [ "/usr/bin/unwise" ]
