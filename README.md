# go-tlog（golang日志组件）
# 如人饮水 冷暖自知
# 1，简介
golang日志组件，类似于log4j;
个人不太喜欢golang当前主流的日志框架类似logrus、zap那种冗余的json格式输出和复杂的代码配置，喜欢java日志框架中的如果配置文件+工厂方法引入+简易的输出方式，所以在写go项目时在没有找到好用的日志框架时选择自己写一个自己喜欢的日志框架，主要借鉴log4j和logback的方式实现的。
# 2，特性
1. 支持配置多种日志输出方式，个性化制定日志header；
2. 支持文件日志最大容量和数量配置；
3. 支持适配其他框架中的日志整合输出（已提供适配gin、gorm框架）；
4. 提供panic方法用于处理error已简化golang中恼人的error处理操作；
# 3，配置文件
在项目中新建config/tlog.yml配置文件：
```
tlog:
     out:  file, console # 输出
     level:
       default: info  # 全局默认级别，info，可以不配置
     format: "%d%t%m [%l] %f"  # 2022/01/01 00:00:00.123 [info] /a/b/c.go:23: 日志输出内容，
     conf:
       - file:
           name: task.log       # log日志名称，默认项目名称
           path: log/           # 支持相对路径和绝对路径
           level: debug         # 会覆盖默认level
           maxSize: 50K         # 日志大小，如果配置了，则会触发文件分割，按照时间划分 eg:100K,10M,1G
           maxNum: 5            # 最大日志保存数，只有设置了maxSize才生效，大于1
       - console:
           level: debug   # 会覆盖默认level
```
# 4，日志级别
支持五种日志输出级别：Trace，Debug，Info，Warn，Error，默认输出级别Info，可自定义指定
# 5，console输出方式
控制台输出方式，Trace，Debug，Info日志级别使用的是stdout（标准输出），Warn，Error日志级别使用的是stderr（标准异常输出）
# 6，file输出方式
文件输出，支持配置文件名称，路径，最大文件大小，最多日志文件数，组件会自动在日志写入到指定大小的文件以时间后缀命名后新建日志文件，并且按照时间自动清理指定数量外的日志文件；
# 7，其他
## 1，适配其他框架
### 1，适配gin框架实现日志输出
gin框架默认是将日志输出到控制台的，我们可以通过一个适配器将gin框架中输出的日志与tLog日志组件适配并完成输出；
```
   # tLogAdapter.go
  var log = tLog.GetTLog()  
   // ------- tLog适配gin ------
   type TLogGinAdapter struct {
   }
   
   func (t *TLogGinAdapter) Write(p []byte) (n int, err error){
   	log.Customize(tLog.Info, "%d%t", string(p[: len(p)-1]))  // tLog提供Customize方法用于定制化的输出，三个参数分别是日志级别，日志头和日志内容
   	return n, err
   }

    # main.go
    func main() {
        // 指定gin的tlog适配器
        gin.DefaultWriter = &tLogAdapter.TLogGinAdapter{}
        gin.DefaultErrorWriter = &tLogAdapter.TLogGinAdapter{}
        // gin.SetMode() gin默认日志级别
    
        // 配置gin
        r := gin.Default()
        // 启动gin
        log.Panic3(r.Run(8080))
    }
```
### 2，适配GORM框架实现日志输出
GORM框架中可以指定自己的日志级别已实现在日志中打印sql等操作，默认输出到控制台的，我们可以使用适配器将日志输出到tLog日志组件中：
```
    # tLogAdapter.go
    var log = tLog.GetTLog()  
    // ------- tLog适配gorm ------
    type TLogGormAdapter struct {
    }
    
    func (t *TLogGormAdapter) Printf(s string, v ...interface{}){
        log.Customize(tLog.Info, "%d%t [%l]", fmt.Sprintf(s, v...))  // tLog提供Customize方法用于定制化的输出，三个参数分别是日志级别，日志头和日志内容
    }
    
    # gormInit.go
    db, err = gorm.Open(mysql.Open(mysqlUrl),
    		&gorm.Config{Logger: logger.New(
    			&tLogAdapter.TLogGormAdapter{},						// 指定输出writer
    			logger.Config{								// 增加配置
    				SlowThreshold: 1 * time.Microsecond,			         // 配置慢sql耗时标准，默认 200 * time.Millisecond
    				LogLevel:      LogLevel,					// 打开Warn级别的日志，其实如果我们不需要修改其他配置比如SlowThreshold可以直接设置当前输出日志级别Warn即可
    		}),
    	})
```

### 3，其他框架可同理实现

## 2，提供Panic方法
golang中大量error异常需要判断并且处理，tLog提供重载的PanicN()方法用于对error的处理并且输出日志：
```
    // 如在gin入参绑定时使用ShouldBind()
	var tasks dto.Tasks
	err := c.ShouldBind(&tasks)
	if err != nil {
		panic(err)  // 当出现异常时需要处理
	}
    
    // 通过tLog改造
	var tasks dto.Tasks
	log.Panic2("入参异常：", c.ShouldBind(&tasks))   // tLog的PanicN方法会判断error的nil来进行日志输出并且panic，可以通过外层统一recover协助
```

```
    tLog提供三种Panic方法
    func (l *TLogger) Panic1(v ...interface{})
    func (l *TLogger) Panic2(s string, err error)
    func (l *TLogger) Panic3(err error)
```
