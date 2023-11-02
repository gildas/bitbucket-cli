FROM golang:1.21 as builder

WORKDIR /

COPY go.mod go.sum /

# If go.mod/go.sum don't change the dependencies will be cached
RUN go mod download

COPY . /

# Build the application
RUN CGO_ENABLED=0 go build -o main .

# ---
FROM alpine:3.18 as system
LABEL org.opencontainers.image.title="bb"
LABEL org.opencontainers.image.description="BitBucket Command Line Interface"
LABEL org.opencontainers.image.authors="Gildas Cherruel <gildas@breizh.org>"
LABEL org.opencontainers.image.licenses="MIT"

# Add CA Certificates and clean
RUN apk update && apk upgrade \
  && apk add ca-certificates \
#  && apk add libcap \
  && rm -rf /var/cache/apk/*

# Creates a harmless user
RUN adduser -D -g '' docker

#set our environment

# Install application, dependencies first
WORKDIR /usr/local/bin
COPY --from=builder /main /usr/local/bin/bb
#RUN setcap 'cap_net_bind_service=+ep' ${APP_ROOT}/cantina

USER docker

# CMD is useless here as this image contains only one binary
ENTRYPOINT [ "/usr/local/bin/bb" ]
