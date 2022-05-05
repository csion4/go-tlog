# go-tlog（golang日志组件）
# 1，简介
golang实现的日志组件，类似于log4j，支持配置多种日志级别和输出方式；
# 2，配置文件
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
# 3，日志级别
支持五种日志输出级别：Trace，Debug，Info，Warn，Error，默认输出级别Info，可自定义指定
# 5，console输出方式
控制台输出方式，Trace，Debug，Info日志级别使用的是stdout（标准输出），Warn，Error日志级别使用的是stderr（标准异常输出）
# 6，file输出方式
文件输出，支持配置文件名称，路径，最大文件大小，最多日志文件数，组件会自动在日志写入到指定大小的文件以时间后缀命名后新建日志文件，并且按照时间自动清理指定数量外的日志文件；


