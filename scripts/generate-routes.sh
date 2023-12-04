#!/bin/sh

dev_routes="src/server/assets/routes"

values="charts/code-editor/values.yaml"
routes="charts/code-editor/routes.yaml"

rm -rf $routes

echo 'routes:' >> $routes
for username in $(yq '.users[].name' $values)
do
  echo "  - name: $username" >> $routes
  echo "    path: $(tr -dc A-Za-z </dev/urandom | head -c 13)" >> $routes
done

# generate routes file for Go
if [ ! -z $1 ]
then
  echo 'dev mode'
  cat $routes > $dev_routes
fi

echo "Routes successfully generated"
