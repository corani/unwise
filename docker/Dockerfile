# Build stage
FROM golang:1.22.3-bookworm AS builder

WORKDIR /src
COPY . . 

ARG CGO_ENABLED=0

RUN mkdir /build                                                    \
 && go build -o /build/usr/bin/unwise ./cmd/unwise/                 \
 && ldd /build/usr/bin/unwise | tr -s '[:blank:]' '\n'              \
  | grep '^/'                                                       \
  | xargs -I % sh -c 'mkdir -p /build/$(dirname %); cp % /build/%;' 

# Final stage
FROM scratch

COPY --from=builder --chown=1000:1000 /build /

USER 1000:1000

EXPOSE 3128

CMD [ "/usr/bin/unwise" ]
