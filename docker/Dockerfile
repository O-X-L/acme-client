FROM --platform=${BUILDPLATFORM} golang:1.25-alpine3.23 AS builder

# BUILD: docker build -f Dockerfile -t acme-client --build-arg "VERSION=1.0" --no-cache --network=host --progress=plain .
# RUN: docker run -i --rm --name acme-client --network=host --volume /etc/acme.yml:/etc/acme.yml --volume /var/lib/acme:/var/lib/acme acme-client:1.0

ARG VERSION=1.0
ARG TARGETPLATFORM
ARG BUILDPLATFORM

ENV SOURCES_URL=https://github.com/O-X-L/acme-client/archive/refs/tags/${VERSION}.tar.gz

WORKDIR $GOPATH/src/github.com/O-X-L/acme-client/src

RUN set -eux; \
        wget -O acme-client.tar.gz "$SOURCES_URL"; \
        tar -xzf acme-client.tar.gz -C $GOPATH/src/github.com/O-X-L/acme-client --strip-components=1; \
        rm acme-client.tar.gz; \
        case ${TARGETPLATFORM} in \
            "linux/amd64")  GOARCH=amd64  ;; \
            # arm64 and arm64v8 are equivalent in go and do not require a goarm
            # https://github.com/golang/go/wiki/GoArm
            "linux/arm64" | "linux/arm/v8")  GOARCH=arm64  ;; \
            "linux/ppc64le")  GOARCH=ppc64le  ;; \
            "linux/arm/v6") GOARCH=arm GOARM=6  ;; \
            "linux/arm/v7") GOARCH=arm GOARM=7 ;; \
        esac; \
        printf "Building acme-client for arch ${GOARCH}\n"; \
        GOARCH=${GOARCH} go build -a -o /acme-client $GOPATH/src/github.com/O-X-L/acme-client/src/cmd/main.go

FROM alpine:3.23

ARG UID=9799
ARG USER=acme

RUN set -eux; \
        addgroup --gid $UID --system $USER; \
        adduser \
                --disabled-password \
                --ingroup $USER \
                --no-create-home \
                --system \
                --uid $UID \
                $USER \
        ; \
        mkdir /app; \
        chown $USER:$USER /app

COPY --from=builder /acme-client /usr/local/bin

USER $USER
WORKDIR /app

CMD ["/usr/local/bin/acme-client", "-path-cnf=/etc/acme.yml"]
