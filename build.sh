#!/bin/bash 

prog=$(realpath "$0")
root=$(dirname "$prog")
rc=0

source "${root}/scripts/helpers.sh"

function print_usage {
  echo "Usage: $prog [options...]"
  echo
  echo "Available options:"
  echo "  -h                    prints this help"
  echo "  -b                    build service"
  echo "  -d                    build the docker image"
  echo "  -g                    generate mocks"
  echo "  -t                    run unit tests"
}

function go_install_tools {
    do_echo tools/install.sh
}

function require_tool {
    export PATH=/go/bin:"${root}/tools/bin":${PATH}

    command -v "${1}" > /dev/null
    if [ $? -eq 1 ]; then
        log_info "${1} not found, installing"
        go_install_tools
    fi

    log_notice "using $(command -v "$1") version $(go version -m "$(command -v "$1")" | awk '$1 == "mod" { print $3 }')"
}

function update_hash {
    VERSION=${VERSION:-dev}

    echo "${VERSION}" > cfg/VERSION

    if [ -z "${GITHUB_SHA:-}" ]; then
        if command -v git > /dev/null; then
            hash=$(git rev-parse HEAD)
        else
            log_warning "git not found, couldn't update hash (set GITHUB_SHA)"
        fi
    else
        hash=${GITHUB_SHA}
    fi

    echo "${hash}" > cfg/HASH
}

function go_build {
    if [ -z "${BUILD_DEBUG:-}" ]; then 
        build="release"
        opts=(-trimpath -ldflags '-s -w')
    else
        build="debug"
        opts=(-gcflags 'all=-N -l')
    fi

    log_info "Building unwise version $(cat cfg/VERSION)/$(cat cfg/HASH) (${build} build)"

    do_echo go build "${opts[@]}"   \
        -o bin/unwise               \
        ./cmd/unwise/
}

function go_test {
    if [ "$(go env GOOS)" == "android" ]; then
        # NOTE(daniel): `-race` is not supported on Android. This check is probably not
        # exhaustive, but it's good enough for now.
        log_warning "skipping race detection on Android"
        opts=()
    else 
        opts=(-race)
    fi

    do_echo go test "${opts[@]}"    \
        -coverprofile=coverage.txt  \
        -covermode=atomic           \
        ./... 
}

function go_generate {
    require_tool mockery 
    require_tool gofumpt 

    do_echo mockery
    do_echo ./scripts/gofmt.sh
}

function docker_build {
    VERSION=${VERSION:-dev}

    do_echo docker build                            \
        --build-arg VERSION="${VERSION}"            \
        --build-arg BUILD_DEBUG="${BUILD_DEBUG}"    \
        -t unwise:"${VERSION}"                      \
        -f docker/Dockerfile                        \
        .
}

mkdir -p bin data

update_hash

if [ "$#" -eq "0" ]; then 
    print_usage
    exit 1
fi

while [ "$#" -gt "0" ]; do
  arg=$1
  shift

  case $arg in
    -h)
        print_usage
        exit 0
        ;;
    -b)
        go_build
        ;;
    -d)
        docker_build
        ;;
    -t)
        go_build
        go_test
        ;;
    -g)
        go_generate 
        ;;
    *)
        log_error "unrecognized argument '$arg'"
        print_usage
        exit 1
        ;;
  esac
done

exit "${rc}"
