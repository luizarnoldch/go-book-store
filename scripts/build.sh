#!/bin/bash

# Obtener el tiempo de inicio
start_time=$(date +%s)

FOLDERS_LAMBDAS=($(ls -d lambdas/*/))

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

build_lambdas() {
    for folder in "${FOLDERS_LAMBDAS[@]}"; do
    (
        folder_name=$(basename "${folder}")
        cd "lambdas/$folder_name" || { echo "Failed to cd into lambdas/$folder_name"; exit 1; }
        go build -tags lambda.norpc -o bootstrap || { echo "Failed to build in lambdas/$folder_name"; exit 1; }
        zip ../../bin/${folder_name}.zip bootstrap || { echo "Failed to zip bootstrap in lambdas/$folder_name"; exit 1; }
        rm -rf bootstrap
    )
    done
}

build_lambdas

# Obtener el tiempo de finalización
end_time=$(date +%s)

# Calcular y mostrar la duración
duration=$((end_time - start_time))
echo "Total build time: $duration seconds"