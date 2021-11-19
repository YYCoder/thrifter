# thrifter
**Non-destructive** parser/printer for [thrift](https://thrift.apache.org/docs/types.html) with zero third-party dependency.

[![YYCoder](https://circleci.com/gh/YYCoder/thrifter.svg?style=svg)](https://app.circleci.com/pipelines/github/YYCoder/thrifter)
[![goreportcard](https://goreportcard.com/badge/github.com/yycoder/thrifter)](https://goreportcard.com/report/github.com/yycoder/thrifter)
[![GoDoc](https://pkg.go.dev/badge/github.com/YYCoder/thrifter)](https://pkg.go.dev/github.com/YYCoder/thrifter)
[![Codecov](https://codecov.io/gh/YYCoder/thrifter/branch/master/graph/badge.svg)](https://codecov.io/gh/YYCoder/thrifter)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

[中文文档](./docs/cn.md)

## Inspiration
There are several thrift parsers on github, but each of them have issues on preserve original format, since mostly, they are used to generate rpc code. But, this project has different purpose, which is focus on helping you write thrift code more efficiently, thanks to its non-destructive code transformation. Since it's a non-destructive parser, we can do a lot of stuff on top of it, such as **code formatting**, **code transforming**, etc.

Currently, it's mainly used by my other little project called [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift), which is a code transformer between protobuf and thrift.

Here are some other thrift parsers on github I discovered before start protobuf-thrift, none of them is 100 percent suitable for it.

### Similar Packages
1. [go-thrift](https://github.com/samuel/go-thrift): mainly used to generate rpc code, ignore white space, comments and lose statements order

2. [thriftrw-go](https://github.com/thriftrw/thriftrw-go): thrift parser and code generator which open sourced by Uber, same issue as above

3. [thriftgo](https://github.com/cloudwego/thriftgo): another thrift parser and code generator, same issue as above

4. [thrift-parser](https://github.com/creditkarma/thrift-parser): thrift parser written in typescript, same issue

So, that's why I started thinking of writing a new thrift parser which preserves all the format.

Thanks to [rocambole](https://github.com/millermedeiros/rocambole), behind which idea is perfect for this project.

### Core Concept
The main idea behind thrifter on achieve non-destructive is that, we use a linked-list to chain all the tokens.

Lets think of the essence of source code, it's just a chain of token, different token combination declare different syntax, so if we want to preserve original format, we must preserve all tokens from the source code.

The best data structure for this chain of tokens is linked-list, since it's easier to modify than array, we only need to change some pointer, and we can patch start token and end token to each ast node, so that we are able to easily iterate over tokens within a node.

When iterate over token linked-list, we also provide a Map structure for each ContainerType, such as enum/struct/service, in order to find the field node started by the token.

## Usage
Initialize Parser first, and specify the io.Reader which consume source code:

```go
parser := thrifter.NewParser(strings.NewReader(XXX), false)
// or
file, err := os.Open(XXX)
if err != nil {
   return nil, err
}
defer file.Close()
parser := thrifter.NewParser(file, false)
```

then, simply use `parser.Parse` to start parsing:

```go
definition, err := parser.Parse(YOUR_FILE_NAME)
```

and that's it, now we have the root node for the source code, which structure like this:

```go
type Thrift struct {
	NodeCommonField
	// thrift file name, if it exists
	FileName string
	// since Thrift is the root node, we need a property to access its children
	Nodes []Node
}
```

You might wonder what the hell is `NodeCommonField` nested into Thrift node, that's the magic of thrifter, we will discuss it in the **AST Node** section.

### Code Print
The most amazing thing about thrifter is that, it is also a non-destructive code printer.

Think about this case, When you want to write a code generator to optimize your workflow, normally you would use a code parser to get the code ast, and then manipulate it. Under some circumstances, you merely want to add some new code to it and leave the rest intact, normal code parser could not able to do that, since they will ignore whitespace like line-breaks/indents.

With thrifter, you can just initialize your new code ast node, and then patch the `StartToken` to the original ast's token linked-list, all other code is unchanged, like this:

```go
// 1. initialize new node, enum, for instance
// for simplicity, you can just initialize a parser to parse the code you want to generate, in order to get the code tokens linked-list
p := thrifter.NewParser(`enum a {
    A = 1
    B = 2;
    C
    D;
}`, false)
startTok := parser.next() // consume enum token
enumNode := NewEnum(startTok, nil)
if err := enumNode.parse(p); err != nil {
   t.Errorf("unexpected error: %v", err)
   return
}

// 2. patch the generated code StartToken to any where you want to put it
preNext := someNodeFromOriginalCode.EndToken.Next
someNodeFromOriginalCode.EndToken.Next = enumNode.StartToken
enumNode.EndToken.Next = preNext

// 3. last, use node.String to print the code
fmt.Println(thriftNodeFromOriginalCode.String())
```

Each `thrifter.Node` have their own `String` function, so, you can also print the node standalone not the whole thrift file.

The principle of `String` is pretty simple, it just traverse the token and write them one by one:

```go
func toString(start *Token, end *Token) string {
	var res bytes.Buffer
	curr := start
	for curr != end {
		res.WriteString(curr.Raw)
		curr = curr.Next
	}
	res.WriteString(end.Raw)
	return res.String()
}
```

Note that, when you manipulate original ast like above, the original `Thrift.Nodes` fields is unchanged, but it doesn't affect code print, since it only iterate over tokens, not nodes. However, you can manually add the node to `Thrift.Nodes` by yourself for consistency.

### AST Node
To understand the idea behind thrifter, there are two struct and one interface you must know:

```go
type NodeCommonField struct {
	Parent     Node
	Next       Node
	Prev       Node
	StartToken *Token
	EndToken   *Token
}

type Token struct {
	Type  token
	Raw   string // tokens raw value, e.g. comments contain prefix, like // or /* or #; strings contain ' or "
	Value string // tokens transformed value
	Next  *Token
	Prev  *Token
	Pos   scanner.Position
}

type Node interface {
	// recursively output current node and its children
	String() string
	// recursively parse current node and its children
	parse(p *Parser) error
	// get node value
	NodeValue() interface{}
	// get node type, value specified from each node
	NodeType() string
}
```

Firstly, `NodeCommonField` is the basic of achieving non-destructive, it will be nested into each ast node, whatever `node.Type` is. These two fields are essential:

* **StartToken**: the start token of the node, which means you can easily iterate over tokens within the node

* **EndToken**: the end token of the node, when iteration within node reaches it, means iteration is done

Second struct `Token` represents a basic token of thrifter, a token can be a symbol, e.g. `-` or `+`, or string literal `"abc"` or `'abc'`, and also a identifier.

> Note that, thrifter considers comment as a token, not a node, currently. I'm not entirely sure it is a good idea, so if some one have questions about it, please open an issue.

And the last interface `Node` represents a thrifter node. Since it's a interface, if you want to access the node fields, you can use `NodeType` to get the type of node, and then do a type assertion of the node:

```go
for _, node := range thrift.Nodes {
    switch 
    case "Namespace":
        n := node.(*thrifter.Namespace)
        fmt.Printf("Namespace: %+v", n)
    case "Enum":
        n := node.(*thrifter.Enum)
        fmt.Printf("Enum: %+v", n)
    case "Struct":
        n := node.(*thrifter.Struct)
        fmt.Printf("Struct: %+v", n)
    case "Service":
        n := node.(*thrifter.Service)
        fmt.Printf("Service: %+v", n)
    case "Include":
        n := node.(*thrifter.Include)
        fmt.Printf("Include: %+v", n)
    }
}
```


## Notice
1. `senum` not supported: since thrift officially don't recommend to use it, thrifter will not handle it, too.

2. current parser implementation is not completely validating `.thrift` definitions, since we think validation feature is better to leave to specific linter.

## Related Packages
Some packages build on top of thrifter:

* [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift): transforming protobuf idl to thrift, and vice versa.

## Contribution
**Working on your first Pull Request?** You can learn how from this *free* series [How to Contribute to an Open Source Project on GitHub](https://kcd.im/pull-request).

### TODO
- [] support comment node
- [] Thrift node support `ElemsMap` to map start token to each element node