#!/bin/bash

# 检查是否提供了目录参数
if [ -z "$1" ]; then
  echo "Usage: $0 <directory>"
  exit 1
fi

# 获取目录路径
DIR="$1"

# 检查目录是否存在
if [ ! -d "$DIR" ]; then
  echo "Directory $DIR does not exist."
  exit 1
fi

# 初始化计数器
counter=1

# 遍历目录中的所有图片文件
for file in "$DIR"/*.{png,jpg,jpeg,gif}; do
  if [ -f "$file" ]; then
    # 获取文件扩展名
    ext="${file##*.}"
    # 构建新的文件名
    new_file="$DIR/$counter.$ext"
    # 重命名文件
    mv "$file" "$new_file"
    # 增加计数器
    ((counter++))
  fi
done

echo "Renamed $((counter-1)) files."