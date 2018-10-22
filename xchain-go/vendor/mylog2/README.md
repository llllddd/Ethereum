# What's Mylog

```go

import "mylog"

var Log ILogger

func main(){
     Log = NewSimpleLogger()
     // Log = NewLogrusLogger()
     // Log = New CologLogger()

     Log.SetLevel(LOG_LEVEL_DEBUG)
     Log.SetPreix("Core")

     Log.Debugln("Hello world !")
     Log.Infoln("Hello mylog !")
}

```
