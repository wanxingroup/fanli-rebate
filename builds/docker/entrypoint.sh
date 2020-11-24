#!/usr/bin/env sh

configFile=${CONFIG_FILE:-/app/config.yaml}
serviceId=${SERVICE_ID:-0}

/app/service --serviceId="${serviceId}" --config="${configFile}"
