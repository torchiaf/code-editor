#!/bin/sh

dir="charts/code-editor/routes"
dev_dir="src/server/routes"

rm -rf "$dev_dir/*.yaml"

# generate users file for Go server
if [ ! -z $1 ]
then
  echo 'dev mode'
  cp -a "$dir/." "$dev_dir/"
fi

echo "Routes files successfully generated"
