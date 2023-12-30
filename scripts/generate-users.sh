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

  i=10000
  for k in $(yq '.users | to_entries | .[].key' $values)
  do
    echo "  - id: id-${i:1}" >> "$dev_dir/users.yaml" 
    echo "    name: $(yq ".users[$k].name" $values)" >> "$dev_dir/users.yaml"
    echo "    password: $(yq ".users[$k].password" $values)" >> "$dev_dir/users.yaml"

    i=$((i+1))
  done

fi

echo "Users file successfully generated"
