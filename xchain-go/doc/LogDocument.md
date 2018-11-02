# 日志mylog2包使用说明

## 1. 初始化日志

将Logger定义为全局变量,并初始化日志,主要定义前缀,日志等级,日志显示格式.详见[3. 日志格式设置](#setlog)

```
prefix := "XChain-Go" //定义前缀

common.InitLog(prefix) // 在common包中初始化Log
common.Logger.SetLevel(mylog.InfoLevel)          //设置显示log等级
common.Logger.SetFlags(mylog.LshortTextNotation) //DEBUG,INFO,ERROR这些字段的显示模式
common.Logger.SetFormatter(mylog.TextFormat)

```
其中common.InitLog()内容如下:
```
var Logger *mylog.SimpleLogger

func InitLog(prefix string) {

	Logger = mylog.NewSimpleLogger(prefix)

}

```

## 2. 调用日志

在其他方法中调用日志时,由于`Logger`为全局变量,并且其为指针,所以初始化的内容可以应用到其他包内的`Logger`变量,举例说明:

```
func PrintLog() {
	Log := common.Logger.NewSessionLogger()
	Log.Debugln("PrintLog hello world!")
	Log.Errorln("PrintLog hello world!")
	Log.Infoln("PrintLog hello world!")
}
```
输出结果:
```
[D] 2018-09-25 15:49:20 <XChain-Go> PrintLog> PrintLog hello world!
[E] 2018-09-25 15:49:20 <XChain-Go> PrintLog> PrintLog hello world!
[I] 2018-09-25 15:49:20 <XChain-Go> PrintLog> PrintLog hello world!
```

调用日志时,也可以不定义:
```
Log := common.Logger.NewSessionLogger()
```
直接使用
```
common.Logger.Debugln("PrintLog hello world!")
````
即:
```
func PrintLog() {
	common.Logger.Debugln("PrintLog hello world!")
	common.Logger.Errorln("PrintLog hello world!")
	common.Logger.Infoln("PrintLog hello world!")
}
```
但是这样输出的就不会打印调用方法,输出结果如下:
```
[D] 2018-09-25 15:49:20 <XChain-Go>PrintLog hello world!
[E] 2018-09-25 15:49:20 <XChain-Go>PrintLog hello world!
[I] 2018-09-25 15:49:20 <XChain-Go>PrintLog hello world!
```

  <span id="setlog"></span>
## 3. 日志格式设置

> 日志等级设置
  
```
eg:common.Logger.SetLevel(mylog.InfoLevel)  
可选日志等级:

FatalLevel   = 1
ErrorLevel   = 3
WarnLevel    = 5
InfoLevel    = 7
DebugLevel   = 9
DDebugLevel  = 11
DDDebugLevel = 13
BackGndLevel = 90
 ```

> 日志格式设置,Flag字段设置

```
eg:common.Logger.SetFlags(mylog.LshortTextNotation) 
可选日志格式:

LshowAttachedInfo       //显示完整日志级别,时间不补0:eg:[DEBUG] 2018-09-25 11:19:4 <XChain-Go>hello world!

LshowAttachedInfoSuffix //显示完整日志级别,时间补全0 eg:[DEBUG] 2018-09-25 11:19:10 <XChain-Go>hello world!

LshortTextNotation      //显示短日志级别 eg:[D] 2018-09-25 11:20:12 <XChain-Go>hello world!

LnoTextNotation         //不显示日志级别 eg: 2018-09-25 11:20:12 <XChain-Go>hello world!

LnoTime                 //不显示时间 eg:[DEBUG] <XChain-Go>hello world!

LnoPrefix               //不显示前缀 eg:[DEBUG] 2018-09-25 11:20:57hello world!
```
> 日志采用Json或文本形式显示

```
common.Logger.SetFormatter(mylog.TextFormat)
可选参数:

JsonFormat  = 1 //以json格式显示 eg:{"timestamp": "2018-09-25 11:22:51","name": "XChain-Go","level": "DEBUG","message": "hello world!"}

TextFormat  = 0 //以文本形式显示 eg:[D] 2018-09-25 11:23:50 <XChain-Go>hello world!

```
