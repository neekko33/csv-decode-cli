# csv-decode-cli

一个使用 Go 编写的交互式 CLI（Bubble Tea），用于将 CSV 中指定字段内的 Unicode 转义（如 `\u65e5\u672c`）转换为可读字符（如 `日本`）。

## 使用方式

```bash
./bin/csv-decode
```

## 编译

```bash
mkdir -p bin
go build -o ./bin/csv-decode .
```

## 交互流程

1. 输入 CSV 文件路径
2. 读取并展示全部表头（header）
3. 手动勾选要转换的字段（可多选）
4. 输入导出 CSV 路径（默认值为输入目录 + `-decoded` 文件名）
5. 如果导出文件已存在，可选择覆盖或返回重新输入

## 按键说明

- `Enter`：确认并进入下一步
- `Space`：字段选择界面勾选/取消勾选
- `Up/Down`：移动光标
- `q` 或 `Ctrl+C`：退出程序

## 帮助

```bash
./bin/csv-decode -h
```