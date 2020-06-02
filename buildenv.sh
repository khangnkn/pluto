#!/usr/bin/env bash
file=$1
prefix=$2

if [[ -z ${file} ]] || [[ -z ${prefix} ]]; then
    echo "Input params not correct"
    echo "Usage: sh './tools/buildenv.sh <file_config_local> <prefix_project>'"
    echo "Example: sh './tools/buildenv.sh ./topaz/configs/admin_tool.yaml zpm_topaz'"
    exit 1
fi

if [[ -e ${file} ]]; then
    echo "Start reading file: $file"
else
    echo "file $file not existed"
    exit 1
fi

echo "\n=========================="

n=1
while read line; do
echo "$line"
n=$((n+1))
done <${file}

echo "\n=== Convert to list ENV with prefix: $prefix ===\n"

s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')

sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
     -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |

awk -F${fs} '{
   indent = length($1)/2;
   vname[indent] = $2;
   for (i in vname) {if (i > indent) {delete vname[i]}}
   if (length($3) > 0) {
      vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
      printf("%s_%s%s: %s\n", toupper("'${prefix}'"),toupper(vn), toupper($2), $3);
   }
}'