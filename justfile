all: binaries

build:
    ./hack/build

shell:
    ./hack/shell

binaries:
    docker buildx bake binaries

binaries-cross:
    docker buildx bake binaries-cross

release:
    ./hack/release $(PLATFORM) $(TARGET)

lint:
    docker buildx bake lint

validate-vendor:
    docker buildx bake validate-vendor

validate-docs:
    docker buildx bake validate-docs

validate-gen:
    docker buildx bake validate-gen

vendor:
    ./hack/update-vendor

docs:
    ./hack/update-docs

outdated:
	docker buildx bake outdated
	cat ./bin/outdated/outdated.txt

gen:
    docker buildx bake update-gen

local:
	docker buildx bake image-default --progress plain

meta:
    docker buildx bake meta  --progress plain

install: binaries
	binname=$(docker buildx bake _common --print | jq -cr '.target._common.args.BIN_NAME') && \
	./bin/build/${binname} install && ${binname} --version

generate: vendor docs gen

validate: lint outdated validate-vendor validate-docs validate-gen

unit-test:
	docker buildx bake unit-test --set "*.args.DESTDIR=/out"
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin:/out unit-test

integration-test:
	docker buildx bake integration-test --set "*.args.DESTDIR=/out"
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock -v ./bin:/out integration-test

test: unit-test integration-test
