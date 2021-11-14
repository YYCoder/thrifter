# thrifter
**Non-destructive** parser for [thrift](https://thrift.apache.org/docs/types.html) with zero third-party dependency.

[![YYCoder](https://circleci.com/gh/YYCoder/thrifter.svg?style=svg)](https://app.circleci.com/pipelines/github/YYCoder/thrifter)
[![GoDoc](https://pkg.go.dev/badge/github.com/YYCoder/thrifter)](https://pkg.go.dev/github.com/YYCoder/thrifter)
[![Codecov](https://codecov.io/gh/YYCoder/thrifter/branch/master/graph/badge.svg)](https://codecov.io/gh/YYCoder/thrifter)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

[中文文档](./docs/cn.md)

## Inspiration
There are several thrift parsers on github, but each of them have issues on preserve original format, since mostly, they are used to generate rpc code. But, this project has different purpose, which is focus on helping you write thrift code more efficiently, thanks to its non-destructive code transformation. Since it's a non-destructive parser, we can do a lot of stuff on top of it, such as **code formatting**, **code transform**, etc.

Currently, it's mainly used by my other little project called [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift), which is a code transformer between protobuf and thrift.

Here are some other thrift parsers on github I discovered before start protobuf-thrift, none of them is 100 percent suitable for it.

### Similar Packages
1. [go-thrift](https://github.com/samuel/go-thrift): mainly used to generate rpc code, ignore white space, comments and lose statements order

2. [thriftrw-go](https://github.com/thriftrw/thriftrw-go): thrift parser and code generator which open sourced by Uber, same issue as above

3. [thriftgo](https://github.com/cloudwego/thriftgo): another thrift parser and code generator, same issue as above

4. [thrift-parser](https://github.com/creditkarma/thrift-parser): thrift parser written in typescript, same issue

So, that's why I started thinking of writing a new thrift parser which preserves all the format.

Thanks to [rocambole](https://github.com/millermedeiros/rocambole), behind which idea is perfect for this project.

## Usage

### AST Node

## Notice
senum not supported.

Current parser implementation is not completely validating .thrift definitions.

## Related Packages
on top of thrifer.

## Contribution
**Working on your first Pull Request?** You can learn how from this *free* series [How to Contribute to an Open Source Project on GitHub](https://kcd.im/pull-request).