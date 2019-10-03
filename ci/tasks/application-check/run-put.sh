#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
base=$PWD

host = "cf apps | grep app-ruby-sample | awk '{print $6}'"
end-point = "http://{$host}/service/mongo/test3"
curl -X PUT -H "Content-Type: application/json" -d '{"data":"sometest130"}' ${end-point}
