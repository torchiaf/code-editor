#!/bin/sh

dev_users="src/server/assets/users"
values="charts/code-editor/values.yaml"

# generate users file for Go
if [ ! -z $1 ]
then
  echo 'dev mode'
  echo "users: " > $dev_users
  yq '.users' $values >> $dev_users
fi

echo "Users file successfully generated"
