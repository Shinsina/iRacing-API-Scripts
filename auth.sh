#!/usr/bin/env bash

set -a
source .env
set +a

go run auth.go
