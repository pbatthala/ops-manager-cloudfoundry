#!/bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

for module in ops-manager-cloudfoundry/src/{mongodb-config-agent,mongodb-service-adapter,smoke-tests}; do
    pushd $module

    echo Linting "$module"

    # run linters
    golangci-lint run --modules-download-mode vendor --timeout 15m

    # fix whitespace
    find vendor -type f -exec grep -Iq . {} \; -print0 | xargs -0 sed -i 's/\r$//'

    # get folder hash before running go mod vendor
    OLDHASHES=$(find vendor -type f -print0 | sort -z | xargs -0 sha1sum)

    go mod vendor

    # fix whitespace again (we don't care if it's wrong)
    find vendor -type f -exec grep -Iq . {} \; -print0 | xargs -0 sed -i 's/\r$//'

    # get folder hash after go mod vendor
    NEWHASHES=$(find vendor -type f -print0 | sort -z | xargs -0 sha1sum)

    popd

    if [[ "$OLDHASHES" != "$NEWHASHES" ]]; then
        echo Vendor is out of date!
        LEFT=$(mktemp)
        RIGHT=$(mktemp)
        echo "$OLDHASHES" >"$LEFT"
        echo "$NEWHASHES" >"$RIGHT"
        echo Hash diff:
        diff "$LEFT" "$RIGHT"
        rm -f "$LEFT" "$RIGHT"
        exit 1
    fi
done
