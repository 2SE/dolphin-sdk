### dolphin-sdk使用方法
该sdk为正对dolphin的grpc服务的封装,可以让开发者不用尽量少的理解dolphin的概念，
在目录example中有简单的server和clinet的实现。

### 对于注册服务的约束
1. 所有业务服务需要以"Service"结尾，如"UserService"的"User"将会被解析为Resource
2. 服务注册后只会解析对外暴露的方法（大写字母开头），且入参为proto.Message,出参为（proto.Message,error）
    对于其他方法不会解析为对外接口
3. 版本控制以"_V"开头，正则为"(_V[1-9]+)+((_[0-9]*)|())"，默认无该后缀的默认为"v1"
 ```
 1. GetUser_V1  action：GetUser version:v1
 2. GetUser_V2_123 action：GetUser version:v2.123
 3. GetUser   action:GetUser  version:v1
 其中1和3都存在是，注册该服务会失败，因为这在我们的框架里面属于相同方法
 ```
### 如何使用
参照example


### dolphin提供外部调用的方法
详见 dolphin-sdk/server/setup.go
