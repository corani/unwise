#!/bin/bash
#
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_DIR=$( cd "${SCRIPT_DIR}/.." && pwd )

source "${REPO_DIR}/scripts/helpers.sh"

function should_skip {
    file="$1"

    # NOTE(daniel): skip deleted files
    if [[ ! -f "${file}" ]]; then
        return 0
    fi

    # NOTE(daniel): skip symlinks (or they would be replaced with a copy of the file!)
    if [[ -h "${file}" ]]; then
        return 0
    fi

    # NOTE(daniel): skip "vendor"
    if [[ "${file}" == vendor/* ]]; then
        return 0
    fi
    if [[ "${file}" == */vendor/* ]]; then
        return 0
    fi

    # NOTE(daniel): skip generated files
    if [[ "${file}" == fakes/* ]]; then
        return 0
    fi

    # NOTE(daniel): skip non-go files
    extension="${file##*.}"
    if [[ "${extension}" != "go" ]]; then
        return 0
    fi

    return 1
}

log_info "using $(command -v gofumpt) version: $(go tool gofumpt --version)"

git ls-files | while read -r file; do
    if should_skip "${file}"; then
        continue
    fi

    log_cmd "go tool gofumpt -l -w ${file}"
    output=$(go tool gofumpt -l -w "${file}")
    if [[ "${output}" != "" ]]; then
        log_info "fixed: ${output}"
    fi
done
