# Reproducible builder image
FROM openshift/origin-release:golang-1.10 as build
WORKDIR /go/src/github.com/openshift/machine-api-operator
# This expects that the context passed to the docker build command is
# the machine-api-operator directory.
# e.g. docker build -t <tag> -f <this_Dockerfile> <path_to_machine-api-operator>
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/machine-api-operator -a -ldflags '-extldflags "-static"' github.com/openshift/machine-api-operator/cmd/machine-api-operator
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/nodelink-controller -a -ldflags '-extldflags "-static"' github.com/openshift/machine-api-operator/cmd/nodelink-controller

LABEL io.openshift.release.operator true

# Final container
FROM openshift/origin-base

COPY --from=build /go/src/github.com/openshift/machine-api-operator/bin/machine-api-operator .
COPY --from=build /go/src/github.com/openshift/machine-api-operator/machines machines
COPY --from=build /go/src/github.com/openshift/machine-api-operator/owned-manifests owned-manifests
COPY --from=build /go/src/github.com/openshift/machine-api-operator/install manifests
COPY --from=build /go/src/github.com/openshift/machine-api-operator/bin/nodelink-controller .
