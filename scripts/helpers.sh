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
