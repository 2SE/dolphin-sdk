### dolphin-sdk使用方法
该sdk为正对dolphin的grpc服务的封装,可以让开发者不用尽量少的理解dolphin的概念，
在目录example中有简单的server和clinet的实现。

### 对于注册服务的约束
1. 所有业务服务需要以"Service"结尾，如"UserService"的"User"将会被解析为Resource
2. 服务注册后只会解析对外暴露的方法（大写字母开头），且入参为proto.Message,出参为（proto.Message,error）
    对于其他方法不会解析为对外接口
3. 版本控制以"_V"开头，正则为"(_V[1-9]+)+((_[0-9]*)|())"，默认无该后缀的默认为"v1"
4. 对于Service内部方法支持的入参数最多为2个，出参可以带一个proto实体，必须带error
 ```
 1. GetUser_V1  action：GetUser version:v1
 2. GetUser_V2_123 action：GetUser version:v2.123
 3. GetUser   action:GetUser  version:v1
 其中1和3都存在是，注册该服务会失败，因为这在我们的框架里面属于相同方法
 ```
### 如何使用
[demo演示](./example)

#### 重要方法描述 dolphin-sdk/server/setup.go
```
//优雅停止服务
func Stop() {
}

//启动服务并注册到dolphin
func Start(c *Config, services ...interface{}) {
}

//单独启动服务
func StartGrpcOnly(c *Config, services ...interface{}) {
}

//发送对其他GRPC服务的调用请求
func SendGrpcRequest(path *pb.MethodPath, info *pb.CurrentInfo, message proto.Message) (*pb.ServerComResponse, error) {
}

```
### 微服务开发注意事项
1. 多处功能相似点，以组件的方式开发，如点赞，他横跨多个场景（范式必定统一），随着业务的拓展，被运用的地方会逐步增加，所以以组件的方式开发会对以后的拓展带来便利。
2. 对于从其他grpc服务获取同类信息的需用批量方法而不要简单for循环便利，这个是偷懒人的通病，因为逻辑复杂时一次性查会带来复杂度，这个对于多服务互相依赖的微服务环境，这是很可怕的。因为你不知道你循环调了5次方法，而提供这个方法的微服务内部也是5个20长度的循环，制造高并发就很容易了（这个虽然不是难点，但很重要，微服务大忌）。
3. 微服务间调用互相交流好，避免装饰模式下的a => b => a 这种场景（a从b服务中取得了信息填充在实体X中，b服务再从a服务中拿到的X并填充其他属性，有些复杂信息会出现这种场景，需要人为避免，因为这在跨编译下是被允许的）。
4. 避免链路过长的调用，会增加无畏的api调用时间
