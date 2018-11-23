#!/bin/bash -eu

source /build-common.sh

BINARY_NAME="dupfinder"
COMPILE_IN_DIRECTORY="cmd/dupfinder"
BINTRAY_PROJECT="joonas/dupfinder"

INCLUDE_WINDOWS="true"

standardBuildProcess
