#!/bin/bash 

prog=$(realpath "$0")
root=$(dirname "$prog")
rc=0

function set_rc {
    if [ "$1" -ne 0 ]; then
        rc=$1
    fi
}

function log_notice {
    label="[INFO]"
    if [[ "${GITHUB_ACTIONS:-false}" == "true" ]]; then
        label="::notice::${label}"
    fi

    # shellcheck disable=SC2145
    echo "${label} $@"
}

function log_info {
    # shellcheck disable=SC2145
    echo "[INFO] $@"
}

function log_warning {
    label="[WARN]"
    if [[ "${GITHUB_ACTIONS:-false}" == "true" ]]; then
        label="::warning::${label}"
    fi

    # shellcheck disable=SC2145
    echo "${label} $@"
}

function log_error {
    label="[ERRO]"
    if [[ "${GITHUB_ACTIONS:-false}" == "true" ]]; then
        label="::error::${label}"
    fi

    # shellcheck disable=SC2145
    echo "${label} $@"
    set_rc 1
}

function log_cmd {
    # shellcheck disable=SC2145
    echo "[CMD ] $@"
}

function do_echo {
    # shellcheck disable=SC2145
    log_cmd "$@"

    if [[ "${GITHUB_ACTIONS:-false}" == "true" ]]; then
        TIMEFORMAT="::debug::[TIME] took %3lR"
    else
        TIMEFORMAT="[TIME] took %3lR"
    fi

    time "$@"
    code=$?
    if [ $code -ne 0 ]; then
        log_error "return code $rc"
        set_rc $code
    fi
}

function print_usage {
  echo "Usage: $prog [options...]"
  echo
  echo "Available options:"
  echo "  -h                    prints this help"
  echo "  -b                    build service"
  echo "  -d                    build the docker image"
  echo "  -t                    run unit tests"
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
    log_info "Building unwise version $(cat cfg/VERSION)/$(cat cfg/HASH)"

    do_echo go build -o bin/unwise ./cmd/unwise/
}

function go_test {
    do_echo go test -race -coverprofile=coverage.txt -covermode=atomic ./... 
}

function docker_build {
    VERSION=${VERSION:-dev}

    do_echo docker build                    \
        --build-arg VERSION="${VERSION}"    \
        -t unwise:"${VERSION}"              \
        -f docker/Dockerfile                \
        .
}

mkdir -p bin data

update_hash

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
    *)
        log_error "unrecognized argument '$arg'"
        print_usage
        exit 1
        ;;
  esac
done

exit "${rc}"
