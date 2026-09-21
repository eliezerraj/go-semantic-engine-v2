#!/usr/bin/env bash

set -euo pipefail

input_file="${1:-phrases2.txt}"
output_file="${2:-phrases2_vector.txt}"
endpoint="http://127.0.0.1:6500/embed"

if [[ ! -f "$input_file" ]]; then
    printf 'Input file not found: %s\n' "$input_file" >&2
    exit 1
fi

temp_file="$(mktemp "${output_file}.XXXXXX")"
trap 'rm -f "$temp_file"' EXIT

while IFS= read -r line || [[ -n "$line" ]]; do
    payload="$(jq -cn --arg input "$line" '{inputs: $input, normalize: true, truncate: false, truncation_direction: "Right"}')"
    response="$(curl --fail-with-body --silent --show-error \
        --request POST \
        --header 'Content-Type: application/json' \
        --data "$payload" \
        "$endpoint")"

    jq -c 'if (type == "array" and length == 1 and (.[0] | type == "array")) then .[0] else . end' <<<"$response" >>"$temp_file"
done < "$input_file"

mv "$temp_file" "$output_file"
trap - EXIT

printf 'Wrote vectors to %s\n' "$output_file"