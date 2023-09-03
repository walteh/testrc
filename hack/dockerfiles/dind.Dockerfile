


FROM gobase AS gotestsum
ARG GOTESTSUM_VERSION
ENV GOFLAGS=
RUN --mount=target=/root/.cache,type=cache <<EOT
	GOBIN=/out/ go install "gotest.tools/gotestsum@${GOTESTSUM_VERSION}" &&
	/out/gotestsum --version
EOT


FROM gobase AS docker
ARG TARGETPLATFORM
ARG DOCKER_VERSION
WORKDIR /opt/docker
RUN <<EOT
CASE=${TARGETPLATFORM:-linux/amd64}
DOCKER_ARCH=$(
	case ${CASE} in
	"linux/amd64") echo "x86_64" ;;
	"linux/arm/v6") echo "armel" ;;
	"linux/arm/v7") echo "armhf" ;;
	"linux/arm64/v8") echo "aarch64" ;;
	"linux/arm64") echo "aarch64" ;;
	"linux/ppc64le") echo "ppc64le" ;;
	"linux/s390x") echo "s390x" ;;
	*) echo "" ;; esac
)
echo "DOCKER_ARCH=$DOCKER_ARCH" &&
wget -qO- "https://download.docker.com/linux/static/stable/${DOCKER_ARCH}/docker-${DOCKER_VERSION}.tgz" | tar xvz --strip 1
EOT
RUN ./dockerd --version && ./containerd --version && ./ctr --version && ./runc --version

FROM gobase AS integration-test-base
ARG BIN_NAME
# https://github.com/docker/docker/blob/master/project/PACKAGERS.md#runtime-dependencies
RUN apk add --no-cache \
	btrfs-progs \
	e2fsprogs \
	e2fsprogs-extra \
	ip6tables \
	iptables \
	openssl \
	shadow-uidmap \
	pigz \
	xfsprogs \
	xz
COPY --link --from=gotestsum /out/gotestsum /usr/bin/
COPY --link --from=registry /bin/registry /usr/bin/
COPY --link --from=docker /opt/docker/* /usr/bin/
COPY --link --from=buildkit /usr/bin/buildkitd /usr/bin/
COPY --link --from=buildkit /usr/bin/buildctl /usr/bin/
COPY --link --from=binaries /${BIN_NAME} /usr/bin/
COPY --link --from=buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
COPY --link  --from=compose-bin /docker-compose /usr/libexec/docker/cli-plugins/docker-compose
