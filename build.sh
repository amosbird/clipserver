#!/usr/bin/env sh

go build -ldflags="-linkmode external -extldflags '-static'" -o /tmp/gentoo/opt/bin/clipserver
