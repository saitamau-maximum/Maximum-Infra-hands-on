#!/bin/bash

# ディレクトリ内のすべてのディレクトリとファイルを再帰的に探す
find "$1" -depth -name '*_*' | while read -r old_name; do
  # スネークケース（例: my_directory_name）をキャメルケース（例: myDirectoryName）に変換
  new_name=$(echo "$old_name" | sed -r 's/(^|_)([a-z])/\U\2/g' | sed 's/_//g')

  # 変更が必要な場合に名前を変更
  if [[ "$old_name" != "$new_name" ]]; then
    mv "$old_name" "$new_name"
    echo "Renamed: $old_name -> $new_name"
  fi
done
