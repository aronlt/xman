#!/bin/bash

# 检查是否在 Git 仓库中
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
  echo "错误：当前目录不是一个 Git 仓库，请切换到 Git 仓库目录后再运行此脚本。"
  exit 1
fi


# 查看未暂存的变更内容
echo ""
echo "------------------------------------------"
echo "1. 工作区中未暂存的文件变更内容 (git diff)"
echo "------------------------------------------"
git diff
if [ $? -ne 0 ]; then
    echo "工作区中未暂存的文件变更无法获取。"
fi

# 查看暂存区中的变更内容
echo ""
echo "------------------------------------------"
echo "2. 暂存区已暂存但未提交的文件变更内容 (git diff --cached)"
echo "------------------------------------------"
git diff --cached
if [ $? -ne 0 ]; then
    echo "暂存区中变更的文件无法获取。"
fi

# 查看文件状态总览
echo ""
echo "------------------------------------------"
echo "3. 所有变更的文件状态清单 (git status)"
echo "------------------------------------------"
git status
if [ $? -ne 0 ]; then
    echo "文件状态无法获取。"
fi