#!/usr/bin/env bash

set -e

BASE_PATH=$PWD

SWAGGER_CLIENT_FOLDER=frontend/src/swagger-client

docker run --rm -v ${PWD}:/local -u `id -u`:`id -g` swaggerapi/swagger-codegen-cli-v3 generate -i /local/docs/swagger.yaml -l typescript-axios -o /local/$SWAGGER_CLIENT_FOLDER

cd $SWAGGER_CLIENT_FOLDER
rm -rf .swagger-codegen
rm .gitignore .swagger-codegen-ignore .npmignore git_push.sh package.json README.md tsconfig.json

for file in apis/*.ts
do
    echo -e "// @ts-nocheck\n$(cat $file)" > $file
done

cd $BASE_PATH
