#!/bin/sh

values="helm-charts/code-editor/values.yaml"
dev_dir="src/server/assets/users"

# Add users in dev-mode Go server
if [ ! -z $1 ]
then
  echo 'dev mode'

  rm -rf $dev_dir
  mkdir $dev_dir
  echo "users: " > "$dev_dir/users.yaml"
  yq '.users' $values >> "$dev_dir/users.yaml"
fi

echo "Users file successfully generated"
