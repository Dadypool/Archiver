#!/bin/bash

files=("bib" "book1" "book2" "geo" "news" "obj1" "obj2" "paper1" "paper2" "pic" "progc" "progl" "progp" "trans")

bin_dir="./bin"
dataset_dir="./dataset"

for file in "${files[@]}"; do
    input_file="${dataset_dir}/${file}"
    intermediate_file="${dataset_dir}/${file}_encoded"
    decoded_file="${dataset_dir}/${file}_decoded"

        echo "Encoding $input_file..."
        $bin_dir/encoder "$input_file" "$intermediate_file"

    if [[ $? -ne 0 ]]; then
        echo "Error encoding file: $file"
        continue
    fi

    echo "Decoding to $decoded_file..."
    $bin_dir/decoder "$intermediate_file" "$decoded_file"

    if [[ $? -ne 0 ]]; then
        echo "Error decoding file: $file"
        continue
    fi

    echo "Finished processing $file"
done
