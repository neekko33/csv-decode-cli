# csv-decode-cli

一个使用 Go 编写的命令行工具，用于将 CSV 中指定字段内的 Unicode 转义（如 `\u65e5\u672c`）转换为可读字符（如 `日本`）。

## 使用方式

```bash
./csv-decode-cli <input.csv> <output.csv> <field1> [field2 field3 ...]
```

## 编译

```bash
go build -o csv-decode-cli .
```