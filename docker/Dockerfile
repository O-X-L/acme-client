FROM golang:alpine AS builder

# BUILD: docker build -f Dockerfile -t acme-client:v1.0.2 --no-cache --network=host --progress=plain ..
# RUN: docker run -i --rm --name acme-client --network=host --volume /etc/acme.yml:/etc/acme.yml --volume /var/lib/acme:/var/lib/acme acme-client

COPY . /tmp/repo

RUN cd /tmp/repo && \
    go mod tidy && \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /usr/local/bin/acme-client ./cmd/main.go

FROM alpine:3.23

ARG UID=9800
ARG USER=acme
ARG ADDITIONAL_GID=113
ARG ADDITIONAL_GROUP=ssl-cert

RUN set -eux; \
        addgroup --gid $UID --system $USER; \
        addgroup --gid $ADDITIONAL_GID --system $ADDITIONAL_GROUP; \
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

COPY --from=builder /usr/local/bin/acme-client /usr/local/bin

USER $USER
WORKDIR /app

CMD ["/usr/local/bin/acme-client", "-c=/etc/acme.yml"]
