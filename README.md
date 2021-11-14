# thrifter
Non-destructive parser for [thrift](https://thrift.apache.org/docs/types.html) with zero third-party dependency.

[![YYCoder](https://circleci.com/gh/YYCoder/thrifter.svg?style=svg)](https://app.circleci.com/pipelines/github/YYCoder/thrifter)
[![GoDoc](https://pkg.go.dev/badge/github.com/YYCoder/thrifter)](https://pkg.go.dev/github.com/YYCoder/thrifter)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)


[中文文档](./docs/cn.md)

## Inspiration

### Similar Packages
1. [go-thrift](https://github.com/samuel/go-thrift)：目前 pb-thrift cli 使用的 thrift 解析器和语法树。主要用来生成 rpc 代码的，问题是会丢失注释、顺序、换行等内容，不适合作为代码格式化工具
2. [thriftrw-go](https://github.com/thriftrw/thriftrw-go)：uber 开源的 thrift ⇒ go 的代码生成器，内置 thrift 的语法解析。问题也是会丢失格式
3. [thriftgo](https://github.com/cloudwego/thriftgo)：也是一个 thrift ⇒ go 的代码生成器，但未来目标是生成多种语言的代码。问题同上
4. [thrift-parser](https://github.com/creditkarma/thrift-parser)：ts 编写的 thrift 解析器

## Usage

### AST Node

## Notice
senum not supported.

Current parser implementation is not completely validating .thrift definitions.

## Related Packages
on top of thrifer.

## Contribution
**Working on your first Pull Request?** You can learn how from this *free* series [How to Contribute to an Open Source Project on GitHub](https://kcd.im/pull-request).