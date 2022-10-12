#!/bin/bash

DBNAME=$1

SHELL_FOLDER=$(cd "$(dirname "$0")" || exit;pwd)
echo "Current wd:${SHELL_FOLDER}"

function remove_old_files() {
#    rm -vf "$SHELL_FOLDER/$1"/*.go
#    find "$SHELL_FOLDER/$1"/* \
#    -type d ! -name 'schema' \
#    -o ! -name 'template' \
#    -print0 | xargs -0 rm -rvf
    return
}

function generate_schema() {
    go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/lock \
        --template "$SHELL_FOLDER/$1"/template \
        --template glob="$SHELL_FOLDER/$1/template/*.tmpl" "$SHELL_FOLDER/$1"/schema
    return
}

function generate(){
    dbname=${DBNAME}
    if [ ! -d "${dbname}" ]; then
      echo "Directory not found: ${dbname}, skipping ${dbname} generation!"
      return
    fi

    echo "Process $dbname module"

    echo " ->delete old $dbname schema file:"
    remove_old_files "$dbname"

    echo " ->generate new $dbname schema:"
    generate_schema "$dbname"
}

generate