#!/bin/bash
# 这是一个辅助脚本，用于自动生成 mock 代码。
# 它会读取我们刚刚创建的接口文件，然后生成一个 mock 文件，供测试时使用。

# 安装 mockgen 工具 (如果尚未安装)
go install go.uber.org/mock/mockgen@latest

echo "正在根据 repository/interfaces.go 生成 mock 文件..."

# -source: 指定接口定义文件的路径 (我们的“图纸”)
# -destination: 指定生成的 mock 文件输出路径
# -package: 指定生成的 mock 文件的包名
mockgen -source=repository/interfaces.go -destination=repository/mocks/mock_repo.go -package=mocks

echo "Mock 文件已成功生成在 repository/mocks/mock_repo.go"