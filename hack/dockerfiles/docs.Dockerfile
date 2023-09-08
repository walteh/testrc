# syntax=docker/dockerfile:labs

ARG GO_VERSION=
ARG BUILDRC_VERSION=

FROM walteh/buildrc:${BUILDRC_VERSION} AS buildrc

FROM alpine AS tools
COPY --from=buildrc /usr/bin/buildrc /usr/bin/buildrc
RUN apk add --no-cache git rsync


FROM golang:${GO_VERSION}-alpine AS docsgen
WORKDIR /src
RUN --mount=target=. \
	--mount=target=/root/.cache,type=cache \
	go build -mod=vendor -o /out/docsgen ./docs/generate.go

FROM tools AS gen
RUN apk add --no-cache rsync git
WORKDIR /src
COPY --from=docsgen /out/docsgen /usr/bin
ARG FORMATS
ARG BUILDX_EXPERIMENTAL
RUN --mount=target=/context \
	--mount=target=.,type=tmpfs <<EOT
	set -e
	rsync -a /context/. .
	docsgen "docs/reference"
	mkdir /out
	cp -r docs/reference/* /out
EOT

FROM scratch AS update
COPY --from=gen /out /

FROM tools AS validate
ARG DESTDIR
COPY ${DESTDIR} /old
COPY --from=update / /new
RUN set -e && buildrc diff --current=/old --correct=/new --glob=**/*
