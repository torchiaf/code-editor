#!/bin/sh

dev_users="src/server/assets/users"

values="charts/code-editor/values.yaml"
users="charts/code-editor/users.yaml"

rm -rf $users

echo 'users:' >> $users
for k in $(yq '.users | to_entries | .[].key' $values)
do 
  echo "  - name: $(yq ".users[$k].name" $values)" >> $users
  echo "    password: $(yq ".users[$k].password" $values)" >> $users
  echo "    path: $(tr -dc A-Za-z </dev/urandom | head -c 13)" >> $users
done

# generate users file for Go
if [ ! -z $1 ]
then
  echo 'dev mode'
  cat $users > $dev_users
fi

echo "Users file successfully generated"
