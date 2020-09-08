#!/bin/sh

# Unless explicitly stated otherwise all files in this repository are licensed
# under the Apache License Version 2.0.
# This product includes software developed at Datadog (https://www.datadoghq.com/).
# Copyright 2020 Datadog, Inc.

# Builds Datadogpy layers for lambda functions, using Docker
set -e

# Move into the root directory, so this script can be called from any directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $DIR/../..

LAYER_DIR=".layers"
LAYER_FILE="datadog_extension"
EXTENSION_DIR="./extensions"

rm -rf $LAYER_DIR
rm -rf $EXTENSION_DIR
mkdir $LAYER_DIR
mkdir $EXTENSION_DIR

echo "Building layer"
cd cmd/serverless
GOOS=linux go build  -tags serverless -o ../../$EXTENSION_DIR/datadog-agent
cd ../..
zip -q -r "${LAYER_DIR}/${LAYER_FILE}" $EXTENSION_DIR
rm -rf ./extensions
echo "Done creating archive $LAYER_FILE"
ls $LAYER_DIR | xargs -I _ echo "$LAYER_DIR/_"
