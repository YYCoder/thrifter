# thrifter
零依赖的 **非破坏性** [thrift](https://thrift.apache.org/docs/types.html) 语法解析器、代码生成器。

[![YYCoder](https://circleci.com/gh/YYCoder/thrifter.svg?style=svg)](https://app.circleci.com/pipelines/github/YYCoder/thrifter)
[![goreportcard](https://goreportcard.com/badge/github.com/yycoder/thrifter)](https://goreportcard.com/report/github.com/yycoder/thrifter)
[![GoDoc](https://pkg.go.dev/badge/github.com/YYCoder/thrifter)](https://pkg.go.dev/github.com/YYCoder/thrifter)
[![Codecov](https://codecov.io/gh/YYCoder/thrifter/branch/master/graph/badge.svg)](https://codecov.io/gh/YYCoder/thrifter)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

## 初衷
写这个项目的初衷源自于我的另一个项目 [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift)，一个 protobuf 与 thrift 互转的命令行工具。

调研了网上很多 thrift 解析器，但都存在若干问题，无法完美满足我的需求：

1. [go-thrift](https://github.com/samuel/go-thrift)：主要用来生成 rpc 代码

2. [thriftrw-go](https://github.com/thriftrw/thriftrw-go)：Uber 开源的 thrifter 解析器以及 rpc 代码生成器

3. [thriftgo](https://github.com/cloudwego/thriftgo)：同上

4. [thrift-parser](https://github.com/creditkarma/thrift-parser)：typescript 编写的 thrift 解析器

如上几个项目都存在的问题：

1. **丢失注释**：因为它们都主要专注于 rpc 代码生成，因此注释对它们来说是非必要的

2. **丢失声明顺序**：比如 go-thrift 中会将 enum 中的 values 存到一个 map 中，这样我们就无法保留其原有的顺序。作为一个代码转换工具，肯定是保留原有的顺序是更好的

3. **丢失 whitespaces**：同丢失注释的原因一样，对于 rpc 代码生成器，代码缩进、换行也是不需要的。但对于一个代码转换工具来说，最好还是能够保留原有代码的缩进以及换行

鉴于以上原因，我决定自己编写一个 **非破坏性** 的 thrift 解析器，因此 thrifter 诞生了。

目前 thrifter 主要用于我的 [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift) 项目中，但其实它还可以做很多事情，比如 **代码格式化**、**无损代码转换** 等等。

与它的思想类似的项目还有 [recast](https://github.com/benjamn/recast) 以及 [rocambole](https://github.com/millermedeiros/rocambole)。

感谢 [rocambole](https://github.com/millermedeiros/rocambole)，thrifter 很大程度上借鉴了其思路。

### 核心概念
thrifter 实现非破坏性的核心在于，**使用一个链表保存了所有 token**。

源代码的本质，其实就是一个 token 链，不同的 token 组合实现了不同的语法，因此，如果我们想要实现非破坏性，就必须要保存所有 token。而最好的保存 token 的数据结构，我认为是链表。因为链表修改起来非常简单，只需要改两个指针，而数组则需要移动后面的所有元素，在大数据量的场景下效率低下。

当我们遍历链表时，由于我们拿到的是 token，因此是不知道当前处在哪个 ast 节点上，为此，我们在每个 ContainerType（即 enum/struct/service）上提供了一个 Map 结构，用于通过 StartToken 获取到对应的 Field ast 节点，这样就能快速判断出当前是否处在 ast 节点上。

## 使用方法
首先，初始化 Parser，通过 io.Reader 读取源代码：

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

再使用 `parser.Parse` 开始解析：

```go
definition, err := parser.Parse(YOUR_FILE_NAME)
```

这样我们就有了源代码的根节点（即 `Thrift` 节点），结构大致如下：

```go
type Thrift struct {
	NodeCommonField
	// thrift file name, if it exists
	FileName string
	// since Thrift is the root node, we need a property to access its children
	Nodes []Node
}
```

你可能会问，`NodeCommonField` 嵌套结构体是什么。这就是实现 thrifter 非破坏性的关键，我们会在 **AST Node** 一节再做解释。

### 代码生成
thrifter 最亮眼的功能在于，它不仅是 parser，同时还是 printer。

想想这个场景，当你需要编写代码生成器来优化你的工作流时，通常你需要使用 parser 来解析出 ast，然后再对其进行增删改。如果你只想增加一些新的代码，对于老的代码不想做任何改动，普通的解析器就无法满足了，因为它们会忽略代码格式甚至注释。而 thrifter 正是为这种场景而生。

你可以手动初始化需要添加的新节点，然后将其与任意 ast 节点的 EndToken 串联即可，其余所有部分都是原样保留，没有任何改动。

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

每一个 `thrifter.Node` 都有 `String` 方法，因此你也可以只打印出当前节点，而不是每次都需要打印整个 thrift 文件。

`String` 方法的原理也很简单，就是遍历 token 流然后依次输出：

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

需要注意的是，当你按如上方法修改 ast 时，源 `Thrift.Nodes` 并没有你新增的节点，但这并不影响代码生成，因为 `String` 方法只关注 token。不过，如果你想要保持一致性，也可以手动将新增节点添加到 `Thrift.Nodes` 中去。

### AST Node
要理解 thrifter 的实现思路，有两个结构体和一个 interface 是必须了解的：

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

首先，`NodeCommonField` 是实现非破坏性的基础，它会被内嵌到每一个 ast 节点中。其中两个字段非常重要：

* **StartToken**: 即为当前节点的第一个 token

* **EndToken**: 为当前节点的最后一个 token，当遍历到该 token 时即可认为遍历结束

`Token` 结构体代表一个基础的 thrifter 中的 token，一个 token 可以是任意符号，如 `-` 或 `+`，也可以是字符串字面量，如 `"abc"` or `'abc'`，或者标识符。

> 需要注意的是，thrifter 中注释只作为 token 存在，而不是一个 ast 节点。目前我还不太确定这样做合不合适，如果有任何问题欢迎提 issue。

最后一个 interface `Node` 代表的是一个 thrifter 节点。里面定义了一些所有节点共有的方法，如 parse、String。由于是 interface，我们如果想获取节点内部的字段，可以使用 golang 的类型断言，如下：

```go
for _, node := range thrift.Nodes {
    switch node.NodeType() {
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


## 注意事项
1. `senum` 目前不支持：因为 thrifter 官方不建议使用

2. 目前的实现并不会校验太多语法规则，语法校验最好交给专门的 linter thrifter 只做基本的解析

## 相关库
以为基础 thrifter 构建的应用：

* [protobuf-thrift](https://github.com/YYCoder/protobuf-thrift): protobuf 与 thrift 互转的命令行工具

## Contribution
**Working on your first Pull Request?** You can learn how from this *free* series [How to Contribute to an Open Source Project on GitHub](https://kcd.im/pull-request).

### TODO
- [] 支持注释节点
- [] `Thrift` 节点支持 `ElemsMap`，从而能够通过 `StartToken` 快速映射到对应节点