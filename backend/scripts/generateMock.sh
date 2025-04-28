#!/bin/bash

# ヘルプメッセージ
usage() {
  echo "Usage: $0 <absolute_path>"
  echo
  echo "This script generates mock files using mockgen."
  echo "You can either pass the absolute path as an argument or input it manually."
  echo
  echo "Example:"
  echo "  $0 /home/user/project/backend/internal/domain/repository/user_repository.go"
  exit 1
}

# 引数が渡されていない場合、標準入力を受け付ける
if [ -z "$1" ]; then
  echo "No argument provided. Please enter the absolute path:"
  read -p "Enter the absolute path to the Go interface file: " SOURCE_PATH
else
  SOURCE_PATH=$1
fi

# 引数として絶対パスを取る
if [ -z "$SOURCE_PATH" ]; then
  echo "Error: No path provided."
  exit 1
fi

# SOURCE_PATHが存在するかチェック
if [ ! -f "$SOURCE_PATH" ]; then
  echo "Error: The file '$SOURCE_PATH' does not exist."
  exit 1
fi

# SOURCE_PATHを相対パスに変換
REL_PATH=$(realpath --relative-to=$(pwd) "$SOURCE_PATH")

# 出力先ディレクトリを "mocks" に置き換え
OUTPUT_DIR=$(dirname "$REL_PATH" | sed 's|internal|test/mocks|')

# 元のファイル名を取得し、_mock.go に変換
FILENAME=$(basename "$SOURCE_PATH" .go)
MOCK_FILENAME="${FILENAME}_mock.go"

# 最終的な出力パス
OUTPUT_PATH="${OUTPUT_DIR}/${MOCK_FILENAME}"

# "mocks" フォルダが存在しない場合は作成
mkdir -p "$OUTPUT_DIR"

# mockgen を実行
mockgen -source="$SOURCE_PATH" -destination="$OUTPUT_PATH"

# 結果を出力
echo "✅ Mock generated at: $OUTPUT_PATH"
