#!/bin/bash
stepname="$2"
namespace="$4"
base_image="$6"
worker_images="$8"

read -ra init_containers <<< "$worker_images"

json_array="[]"

for element in "${init_containers[@]}"; do

echo "$element"

commandName="${element%%-*}"

echo "kubectl logs "$base_image" -c "$element" -n "$namespace" "

result=$(kubectl logs "$base_image" -c "$element" -n "$namespace")

patchArray='[{"taskName":"$stepName","commandStatus":[{"commandName":"","commandStatus":""}],"isFail":""}]'

echo  $result >> consolidated_log_$base_image.log

done

echo "Log created successfully"