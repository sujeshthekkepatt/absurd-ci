#!/bin/bash


base_image="$2"

worker_images="$4"

echo "$base_image"
echo "$worker_images"

IFS=","

read -ra init_containers <<< "$worker_images"
json_array="[]"

for element in "${init_containers[@]}"; do
  echo "$element"

echo "kubectl logs "$base_image" -c "$element" -n hola"
result=$(kubectl logs "$base_image" -c "$element" -n hola)

# json_array=$(echo "$json_array" | jq --arg key '$element' --arg value '$result' '. + [{ ($key): $value }]')

    jeResult=$(printf "%s" "$result" | jq -Rs '.')
    echo $jeResult 
  json_object=$(jq -n --arg key "$element" --arg greeting "$jeResult" '{"'"$element"'": $greeting}')

  json_array+=("$json_object")

done
json_string=$(IFS=,; echo "${json_array[*]}")

echo "[$json_string]"


# todo curl json