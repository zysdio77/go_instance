Log - Change log 1.x
==============================

## 1.3

### 功能：

#### 1.【调整】默认日志记录器由zap改回logrus
zap的fields添加函数存在bug。


## 1.2

### 功能：

#### 1.【新增】MyLogger的WithFields方法
允许使用方添加额外字段到JSON格式的日志内容中。

#### 2.【调整】默认日志记录器由logrus改为zap
zap的JSON格式输出更加合理。


## 1.1

### 功能：

#### 1.【新增】uber-go/zap
zap有着与logrus相当的可扩展性，并且JSON输出结构更合理、性能更佳。

#### 2.【调整】日志输出格式统一为JSON
JSON格式的日志内容更利于分析和筛选。
