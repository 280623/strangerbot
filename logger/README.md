# Golang Log Library

基于zap(uber)封装的日志库
使用过程中，应尽量避免显式调用uber zap package的变量与方法，应尽量调用azoya/lib/log的变量与方法

比如 zap.DebugLevel == zapcore.DebugLevel == log.DebugLevel
应当优先使用：log package

## 安装


>   注意：package路径与仓库路径非一致

package path: azoya/lib/log

repo: https://git.coding.net/azoya/azoya_go_lib.git

## 快速开始

    // 使用生产模式配置
    cfg := log.ProductionConfig()
    
    // 修改应用名称
    cfg.Program = "ExampleBasic"
    
    // 生成Logger
    logger, err := cfg.Build()
    if err != nil {
        panic(err)
    }
    
    // 程序退出时刷新缓冲区的日志数据
    defer logger.Sync()
    
    // 创建命名的子记录器
    logger1 := logger.Named("lists")
    logger2 := logger.Named("search")
    
    // 打印日志
    logger.Info("root", log.Int64("int64", 100))
    logger1.Info("lists")
    logger2.Info("search")
    
    // 从context.Context获取opentracing.Span生成 span logger
    // 如果opentracing.Span == nil 那么不会创建子记录器，直接返回原记录器
    spanlogger := listsLogger.SpanFromContext(ctx)
    spanlogger.Info("Info1", log.String("string", "Info"))
    
    // 从opentracing.Span生成 span logger
    // 如果opentracing.Span == nil 那么不会创建子记录器，直接返回原记录器
    spanlogger := listsLogger.Span(span)
    spanlogger.Info("Info1", log.String("string", "Info"))
    
    // 记录器增加或者修改配置项
    logger.Info("hook options test")
    // 增加HOOK配置项
    logger = logger.WithOptions( log.Hooks(func(entry log.Entry) error {
        fmt.Println(entry.Message)
        return nil
    }))
    logger.Info("hook options test")
    
    
[前往更多使用示例](https://coding.net/u/azoya/p/azoya_go_lib/git/blob/feature/upDesign/log/tests/example_test.go)

## 配置信息详解

    // 配置信息
    type Config struct {
        Program           string                 // 应用程序名称
        Level             AtomicLevel            // 记录的最低级别
        Development       bool                   // 是否开发模式
        DisableCaller     bool                   // 是否禁用获取调用者信息
        DisableStacktrace bool                   // 是否禁用获取堆栈跟踪
        DisableSpanLogger bool                   // 是否禁用Tracing的Span Logger
        SpanLevel         AtomicLevel            // Tracing Span Logger 记录的最低级别
        EncoderConfig     EncoderConfig          // 日志编码器配置 推荐使用标准编码器配置
        Sampling          *SamplingConfig        // 日志采样配置 目的是限制住CPU/IO的负载 此处配置的值是每秒
        InitialFields     map[string]interface{} // 记录器自带内置上下文字段
    
        StdoutEnabled   bool             // 是否开启标准输出
        FileEnabled     bool             // 是否开启文件输出
        FileCoreConfig  *FileCoreConfig  // 文件输出配置
    }

### 预置配置

    // 生产环境推荐配置
    // 默认采用文件输出方式 采用内置的标准轮转策略
    func ProductionConfig() *Config {
    
        cfg := &Config{
            Level:          NewAtomicLevelAt(InfoLevel),
            SpanLevel:      NewAtomicLevelAt(InfoLevel),
            Development:    false,
            EncoderConfig:  NewStandardEncoderConfig(),
            FileEnabled:    true,
            FileCoreConfig: StandardFileCoreConfig(),
        }
    
        return cfg
    }

    // 开发环境推荐配置
    // 默认采用文件输出方式 采用内置的标准轮转策略
    func DevelopmentConfig() *Config {
    
        cfg := &Config{
            Level:          NewAtomicLevelAt(DebugLevel),
            SpanLevel:      NewAtomicLevelAt(DebugLevel),
            Development:    true,
            EncoderConfig:  NewStandardEncoderConfig(),
            FileEnabled:    true,
            FileCoreConfig: StandardFileCoreConfig(),
        }
    
        return cfg
    }

### 文件输出配置

    // 文件输出配置
    type FileCoreConfig struct {
        Filename      string // 保存的文件路径
        DisableRotate bool   // 是否禁用轮转策略 另外使用系统脚本去控制轮转
        MaxSize       int    // 单个日志文件最大存储 单位MB
        MaxBackups    int    // 最多保存的文件数
        MaxAge        int    // 最多保存的天数
    }
    
    // 标准化文件输出配置
    func StandardFileCoreConfig() *FileCoreConfig {
    	return &FileCoreConfig{Filename: "./logs/app.log", MaxSize: 50, MaxBackups: 30, MaxAge: 30}
    }
    
    // 修改文件轮转配置
	cfg := log.ProductionConfig()
	cfg.Program = "ExampleFileCoreConfig"
	cfg.FileCoreConfig.Filename = "./logs/lists/app.log" // 保存的文件路径
	cfg.FileCoreConfig.MaxSize = 10                      // 单个日志文件最大存储 单位MB
	cfg.FileCoreConfig.MaxBackups = 15                   // 最多保存的文件数
	cfg.FileCoreConfig.MaxAge = 15                       // 最多保存的天数

	listsLogger, _ := cfg.Build()
	defer logger.Sync()
	
### 采样配置

>   默认配置没有开启采样
>   * 采样是为了降低CPU/IO的负载，对于日志审计不敏感的业务，比如1秒日志打印超过几百条甚至上千可以考虑开启采样
>   * 采样算法：每秒内同一个level与message计数数量如果大于Initial，超出的部分每Thereafter条，只采集1条

	cfg := log.ProductionConfig()
	// 采样配置
	// 按此配置，最后落地的日志数据应该为 100 + （10000 - 100 ) / 100 = 199条
	cfg.Sampling = &log.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	logger, _ := cfg.Build()
	defer logger.Sync()

	for i := 0; i < 10000; i++ {
		logger.Info("sampling test")
	}

    

官方关于采样的说明

Why sample application logs?

Applications often experience runs of errors, either because of a bug or because of a misbehaving user. Logging errors is usually a good idea, but it can easily make this bad situation worse: not only is your application coping with a flood of errors, it's also spending extra CPU cycles and I/O logging those errors. Since writes are typically serialized, logging limits throughput when you need it most.

Sampling fixes this problem by dropping repetitive log entries. Under normal conditions, your application writes out every entry. When similar entries are logged hundreds or thousands of times each second, though, zap begins dropping duplicates to preserve throughput.


### 修改记录最低Level

    // 修改配置内的Level与SpanLevel
    cfg := log.ProductionConfig()
    cfg.Program = "ExampleLevelConfig"
    cfg.Level = log.NewAtomicLevelAt(log.WarnLevel)
    cfg.SpanLevel = log.NewAtomicLevelAt(log.WarnLevel)

    logger, _ := cfg.Build()
    defer logger.Sync()

    logger.Info("test")
    logger.Warn("test")

    // 通过HTTP RESTFUL API修改动态调整
	http.HandleFunc("/logger/level", cfg.Level.ServeHTTP)
	http.HandleFunc("/logger/span/level", cfg.SpanLevel.ServeHTTP)
	if err := http.ListenAndServe("127.0.0.1:9090", nil); err != nil {
		panic(err)
	}
	
## 扩展WriterSync

日常中，可能需要将日志复制写入到更多的地方时，可以通过扩展WriterSync接口方式实现
Encoder 有特殊需求可以自己编写

[前往扩展WriterSync范例](https://coding.net/u/azoya/p/azoya_go_lib/git/blob/feature/upDesign/log/tests/example_test.go#L211)

## 日志字段

| 字段名称 | span logger | 字段说明 |
| ------ | ------ | ------ |
| level | 是 | 日志级别 |
| timestamp | 否 | 时间 ISO8601 |
| msg | 是 | 事件概述 |
| logger | 否 | 记录器名称 默认无名称 |
| program | 否 | 应用程序名称 默认无名称 |
| caller | 否 | 日志调用处 |
| stacktrace | 否 | 调用堆栈信息 |
| 上下文字段 | 是 | 其它描述的上下文信息 |

>   建议同一个上下文字段名称始终保持同一个数据类型


## 日志级别

| 级别 | 级别描述 |
| --- | ---|
| debug | 系统运行中的调试信息，便于开发人员进行错误分析和修正，一般用于程序日志 |
| info | 系统运行的主要关键时点的操作信息，一般用于记录业务日志 |
| warn | 表明会出现潜在错误的情形 |
| error | 指出虽然发生错误事件，但仍然不影响系统的继续运行 |
| fatal | 表示需要立即被处理的系统级错误。当该错误发生时，表示服务已经出现了某种程度的不可用，系统管理员需要立即介入 |

> fatal 打印完应用程序会自动退出

如果一个程序运行不正常，但有充足的DEBUG埋点，那么可以按照下面的操作步骤进行：
>   可以通过HTTP RESTFUL API形式去控制最低记录level
1.  修改日志配置记录级别，开启DEBUG
2.  允许程序一段时间内写入“DEBUG”日志消息
3.  修改日志配置记录级别，关闭DEBUG

## 参考文档

ZAP 相关

*   [zap官方文档](https://godoc.org/go.uber.org/zap)
*   [zap FAQ](https://github.com/uber-go/zap/blob/master/FAQ.md)
*   [深度 | 从Go高性能日志库zap看如何实现高性能Go组件](https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw)

日志 相关

*   [metrics-tracing-and-logging](https://peter.bourgon.org/blog/2017/02/21/metrics-tracing-and-logging.html)
*   [when-to-use-the-different-log-levels](https://stackoverflow.com/questions/2031163/when-to-use-the-different-log-levels)
*   [what-60000-customer-searches-taught-us-about-logging-in-json](https://www.loggly.com/blog/what-60000-customer-searches-taught-us-about-logging-in-json/)
*   [8-handy-tips-consider-logging-json](https://www.loggly.com/blog/8-handy-tips-consider-logging-json/)