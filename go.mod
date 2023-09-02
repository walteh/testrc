module github.com/walteh/testrc

go 1.21.0

replace github.com/moby/buildkit => github.com/moby/buildkit v0.12.1-0.20230831174303-e2dfabdebce2

require (
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/config v1.18.16
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.39
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5
	github.com/docker/buildx v0.11.2
	github.com/docker/distribution v2.8.2+incompatible
	github.com/fatih/color v1.15.0
	github.com/jedib0t/go-pretty/v6 v6.4.7
	github.com/k0kubun/pp/v3 v3.2.0
	github.com/moby/buildkit v0.12.2
	github.com/ory/dockertest/v3 v3.10.0
	github.com/rs/zerolog v1.30.0
	github.com/spf13/afero v1.9.5
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/walteh/buildrc v0.12.7
	github.com/walteh/snake v0.5.0
	golang.org/x/mod v0.12.0
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20230106234847-43070de90fa1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230717121422-5aa5874ade95 // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.15.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.6 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/containerd/containerd v1.7.2 // indirect
	github.com/containerd/continuity v0.4.1 // indirect
	github.com/containerd/typeurl/v2 v2.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v24.0.5+incompatible // indirect
	github.com/docker/docker v24.0.5-0.20230714235725-36e9e796c6fc+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.4.1 // indirect
	github.com/go-git/go-git/v5 v5.8.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/in-toto/in-toto-golang v0.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/opencontainers/runc v1.1.9 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.4.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/shibumi/go-pathspec v1.3.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
