#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x
base=$PWD
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/cf-helper.sh"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
make_pcf_metadata

cf_login
delete_service_app_if_exists "test-mongodb-service" "app-ruby-sample"
cf logout
