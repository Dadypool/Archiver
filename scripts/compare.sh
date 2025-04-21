#!/bin/bash

files=("bib" "book1" "book2" "geo" "news" "obj1" "obj2" "paper1" "paper2" "pic" "progc" "progl" "progp" "trans")

bin_dir="./bin"
dataset_dir="./dataset"

for file in "${files[@]}"; do
    input_file="${dataset_dir}/${file}"
        $bin_dir/compare "$input_file"
done
