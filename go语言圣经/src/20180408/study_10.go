package main


/*
	第十章　包和工具

	Go语言开源包，它们可以通过 http://godoc.org 检索
*/


//一，包简介
//1.包系统设计的目的是为了简化大型程序的设计和维护工作
//2.每个包一般都定义了一个不同的名字空间用于它内部的每个标识符的访问
//3.Go语言三个特性1->包必须在每个文件的开头显式声明
//				2->禁止包的环状依赖
//				3->编译后包的目标文件记录包本身的信息和包的依赖关系


//二，导入路径
//1.每个包是由一个全局唯一的字符串所标识的导入路径定位
//2.Go语言的规范并没有指明包的导入路径字符串的具体含义，导入路径的具体含义是由构建工具来解释的


//三，包声明
//1.每个Go语言源文件的开头都必须有包声明语句
//2.默认包名一般采用导入路径名的最后一段的约定，
// 也有三种例外情况1->main包本身的导入路径是无关紧要的
//				2->文件名是以_test.go为后缀的Go源文件，都由go test命令独立编译
//				3->依赖版本号的管理工具会在导入路径后追加版本号信息，包的名字并不包含版本号后缀

//四，导入声明
//1.每个导入声明可以单独指定一个导入路径
//2.导入的包之间可以通过添加空行来分组
//3.包的导入顺序无关紧要
//4.导入包的重命名：mrand "math/rand"，导入包的重命名只影响当前的源文件，解决名字冲突，简洁包名

//五，包的匿名导入
//1._来重命名导入的包，为空白标识符，并不能被访问，被称为包的匿名导入


//六，包和命名
//1.命名名字要简洁明了
//2.一般采用单数的形式
//3.避免包名有其它的含义
//4.只描述了单一的数据类型，例如html/template和math/rand等

//七，工具

/*$ go
Go is a tool for managing Go source code.

Usage:

	go command [arguments]

The commands are:

	build       compile packages and dependencies
	clean       remove object files
	doc         show documentation for package or symbol
	env         print Go environment information
	bug         start a bug report
	fix         run go tool fix on packages
	fmt         run gofmt on package sources
	generate    generate Go files by processing source
	get         download and install packages and dependencies
	install     compile and install packages and dependencies
	list        list packages
	run         compile and run Go program
	test        test packages
	tool        run specified go tool
	version     print Go version
	vet         run go tool vet on packages

Use "go help [command]" for more information about a command.

Additional help topics:

	c           calling between Go and C
	buildmode   description of build modes
	filetype    file types
	gopath      GOPATH environment variable
	environment environment variables
	importpath  import path syntax
	packages    description of package lists
	testflag    description of testing flags
	testfunc    description of testing functions
*/

//1.为了达到零配置的设计目标，Go语言的工具箱很多地方都依赖各种约定
//2.可以根据导入路径找到存储代码仓库的远程服务器的URL
//3.GOPATH对应的工作区目录有三个子目录1->src子目录用于存储源代码
//								  2->pkg子目录用于保存编译后的包的目标文件
//								  3->bin子目录用于保存编译后的可执行程序
//4.GOROOT用来指定Go的安装目录
//5.-u命令行标志参数，go get命令将确保所有的包和依赖的包的版本都是最新的，并重新编译安装
//6.go build -i命令将安装每个目标所依赖的包
//7.go install命令和go build命令很相似，但是它会保存每个包的编译成果，而不是将它们都丢弃
//8.Go语言的编码风格鼓励为每个包提供良好的文档
//9.go doc命令，该命令打印其后所指定的实体的声明与文档注释
//10.godoc，它提供可以相互交叉引用的HTML页面，但是包含和go doc命令相同以及更多的信息
//11.go list命令可以查询可用包的信息,还可以用"..."表示匹配任意的包的导入路径,其中-json命令行参数表示用JSON格式打印每个包的元信息
//12.go test命令运行Go语言程序中的测试代码

func main() {
	
}
