#!/bin/bash

set -e

(cf uninstall-plugin "FirehosePlugin" || true) && go build -o nozzle main.go && cf install-plugin nozzle
