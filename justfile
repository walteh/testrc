##################################################################
# GENERATE
##################################################################

generate:
	docker buildx bake generate

generate-buf:
	docker buildx bake generate-buf

generate-mockery:
    docker buildx bake generate-mockery

generate-meta:
    docker buildx bake meta

generate-vendor:
    docker buildx bake generate-vendor

generate-docs:
    docker buildx bake generate-docs

##################################################################
# VALIDATE
##################################################################

validate:
	docker buildx bake validate

lint:
    docker buildx bake lint

validate-vendor:
    docker buildx bake validate-vendor

validate-docs:
    docker buildx bake validate-docs

validate-gen:
    docker buildx bake validate-gen

outdated:
	docker buildx bake outdated
	cat ./bin/outdated/outdated.txt

##################################################################
# TEST
##################################################################

test:
	docker buildx bake test
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin/testreports:/out unit-test
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin/testreports:/out integration-test

unit-test:
	docker buildx bake unit-test
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin/testreports:/out unit-test

integration-test:
	docker buildx bake integration-test
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin/testreports:/out integration-test

##################################################################
# BUILD
##################################################################

build:
    docker buildx bake build

release:
    docker buildx bake release

package:
	RELEASE_OUTPUT=$(mktemp -d -t release-XXXXXXXXXX) && \
	docker buildx bake release --set "*.output=${RELEASE_OUTPUT}" && \
	docker buildx bake package --set "*.contexts.released=${RELEASE_OUTPUT}" && \
	rm -rf ${RELEASE_OUTPUT}

local:
	docker buildx bake image-default

install: build
	binname=$(docker buildx bake _common --print | jq -cr '.target._common.args.BIN_NAME') && \
	./bin/build/${binname} install && ${binname} --version
