#!/bin/sh

source_dir="src/templates"
helm_dir="helm-charts/code-editor/assets"
dev_dir="src/server/assets/templates"

# Copy k8s templates to helm-charts
cp -a "$source_dir/." "$helm_dir/"

# Add templates in dev-mode Go server
if [ ! -z $1 ]
then
  echo 'dev mode'

  cp -a "$source_dir/." "$dev_dir/"
fi

echo "Templates files successfully generated"
