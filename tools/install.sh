#!/bin/bash 

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_DIR=$( cd "${SCRIPT_DIR}/.." && pwd )
rc=0

source "${REPO_DIR}/scripts/helpers.sh"

export GOBIN="${SCRIPT_DIR}/bin"
cd "${SCRIPT_DIR}" || exit 1

log_info "found tools:"
tools=$(grep "_" "${SCRIPT_DIR}/tools.go" | cut -d'"' -f2)
for tool in ${tools}; do
    log_info " - ${tool}"
done

for tool in ${tools}; do
    do_echo go install "${tool}"
done

exit "${rc}"
