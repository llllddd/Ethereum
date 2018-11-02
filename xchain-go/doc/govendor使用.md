# govendor
参考链接：
https://blog.csdn.net/yeasy/article/details/65935864

## 1. 安装
> go get -u github.com/kardianos/govendor
 
命令行执行
> govendor

若正常安装，则会输出以下信息
```
govendor (v1.0.9): record dependencies and copy into vendor folder
	-govendor-licenses    Show govendor's licenses.
...

```
## 2. govendor常用的使用命令


> **1.初始化**
进入项目目录，初始化govender。会产生一个包含vender.json的vender目录
```
govendor init
```
> **2.Add packages from $GOPATH**
将gopath中的包添加到vendor中
```
govendor add +external
```
> **3.更新外部所有依赖**
```
govendor update +external
```

## 3. 所有的操作命令
对于 govendor 来说，主要存在三种位置的包：项目自身的包组织为本地（local）包；传统的存放在 $GOPATH 下的依赖包为外部（external）依赖包；被 govendor 管理的放在 vendor 目录下的依赖包则为 vendor 包。
常见的命令如下，格式为 govendor COMMAND。

| 命令      |    功能 | 
| :-------- | :--------|
|init	|初始化 vendor 目录|
|list	|列出所有的依赖包|
|add	|添加包到 vendor 目录，如 govendor add +external 添加所有外部包|
|add PKG_PATH	|添加指定的依赖包到 vendor 目录|
|update	|从 $GOPATH 更新依赖包到 vendor 目录|
|remove	|从 vendor 管理中删除依赖|
|status	|列出所有缺失、过期和修改过的包|
|fetch	|添加或更新包到本地 vendor 目录|
|sync	|本地存在 vendor.json 时候拉去依赖包，匹配所记录的版本|
|get	|类似 go get 目录，拉取依赖包到 vendor 目录|

通过指定包类型，可以过滤仅对指定包进行操作。
具体来看，这些包可能的类型如下：
| 状态      |    缩写状态 | 含义|
| :-------- | :----:| :--------|		
|+local |l	|本地包，即项目自身的包组织|
|+external	|e	|外部包，即被 $GOPATH 管理，但不在 vendor 目录下|
|+vendor	|v	|已被 govendor 管理，即在 vendor 目录下|
|+std	|s	|标准库中的包|
|+unused	|u|	未使用的包，即包在 vendor 目录下，但项目并没有用到|
|+missing|	m	|代码引用了依赖包，但该包并没有找到|
|+program	|p	|主程序包，意味着可以编译为执行文件|
|+outside	 |	|外部包和缺失的包|
|+all	 |	|所有的包|



执行govendor会打印出的使用命令信息
```
govendor (v1.0.9): record dependencies and copy into vendor folder
	-govendor-licenses    Show govendor's licenses.
	-version              Show govendor version
	-cpuprofile 'file'    Writes a CPU profile to 'file' for debugging.
	-memprofile 'file'    Writes a heap profile to 'file' for debugging.

Sub-Commands

	init     Create the "vendor" folder and the "vendor.json" file.
	list     List and filter existing dependencies and packages.
	add      Add packages from $GOPATH.
	update   Update packages from $GOPATH.
	remove   Remove packages from the vendor folder.
	status   Lists any packages missing, out-of-date, or modified locally.
	fetch    Add new or update vendor folder packages from remote repository.
	sync     Pull packages into vendor folder from remote repository with revisions
  	             from vendor.json file.
	migrate  Move packages from a legacy tool to the vendor folder with metadata.
	get      Like "go get" but copies dependencies into a "vendor" folder.
	license  List discovered licenses for the given status or import paths.
	shell    Run a "shell" to make multiple sub-commands more efficient for large
	             projects.

	go tool commands that are wrapped:
	  "+status" package selection may be used with them
	fmt, build, install, clean, test, vet, generate, tool

Status Types

	+local    (l) packages in your project
	+external (e) referenced packages in GOPATH but not in current project
	+vendor   (v) packages in the vendor folder
	+std      (s) packages in the standard library

	+excluded (x) external packages explicitly excluded from vendoring
	+unused   (u) packages in the vendor folder, but unused
	+missing  (m) referenced packages but not found

	+program  (p) package is a main package

	+outside  +external +missing
	+all      +all packages

	Status can be referenced by their initial letters.
...

```