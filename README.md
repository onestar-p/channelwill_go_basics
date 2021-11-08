# golang基础框
> grpc-gateway + gRPC

#### 环境要求
- golang版本：1.17
- consul: V1.10.3

#### 使用扩展
- go.uber.org/zap ：用于打印日志
- entgo.io/ent ：ORM

## 目录/文件说明
```
root
├── config // 配置文件目录
│   └── cert // 私钥、公钥存放目录
│       ├── private.key 
│       └── public.key
│
├── dao // 数据访问代码存放目录
│   └── client
│       └── client.go // 数据库初始化目录（MySQL、Redis...）
│
├── genProto.sh // 生成proto文件
│
├── global // 全局目录
│   └── global.go
│
├── initialize // 初始化文件目录
│   ├── config.go // 配置文件初始化
│   └── logger.go // zap日志初始化
│
├── interceptor // gRPC拦截器目录
│
├── proto // proto文件存放目录
├── service // gRPC 接口目录
└── utils // 工具类目录
```


## 创建项目步骤

### 一、生成proto pb文件
- 涉及文件：root/proto/auth.yaml
- 涉及文件：```root/genProto.sh```

 在生成pb文件时，需要为gRPC接口定义请求路劲；例如：```root/proto/auth.yaml```


说明：

**.yaml
```
type: google.api.Service
config_version: 3

http: 
  rules:
    - selector: ProtoApi.Auth.Login
      post: /v1/login
      body: "*"
    - selector: ProtoApi.Auth.GetUserToken
      get: /v1/get_user_token
```


genProto.sh
```
function genProto {
    PROTO_NAME=$1
    PROTO_PATH=./proto
    if [ $2 ]; then
        GO_OUT_PATH="${PROTO_PATH}/gen/${PROTO_NAME}/${2}"
    else
        GO_OUT_PATH=${PROTO_PATH}/gen/${PROTO_NAME}
    fi
    mkdir -p $GO_OUT_PATH
    protoc -I=${PROTO_PATH} --go_out=plugins=grpc,paths=source_relative:$GO_OUT_PATH ${PROTO_NAME}.proto
    protoc -I=${PROTO_PATH} --grpc-gateway_out=paths=source_relative,grpc_api_configuration=${PROTO_PATH}/${PROTO_NAME}.yaml:$GO_OUT_PATH ${PROTO_NAME}.proto

}

// genProto {proto文件名} + {版本}
genProto etranslate v1
genProto auth v1
```

### 二、Gateway 注册服务
- 涉及文件：```root/utils/register/server/server.go```
- 涉及方法：```func NewClientServices() []*ClientService```
```
func NewClientServices() []*ClientService {
	services := []*ClientService{
		{
			Name:         "auth",
			RegisterFunc: authpb.RegisterAuthServiceHandler, // auth pb
		},
		{
			Name:         "extranslate",                                 // 服务名
			RegisterFunc: etranslatepb.RegisterEtranslateServiceHandler, // etranslate pb
		},
	}
	return services
}
```

### 三、服务注册
- 涉及文件：```root/cmd/server/main.go```

```
...

if err := service.RunGRPCServer(&service.GRPCConfig{
	Name:              appConfig.Name,
	Addr:              addr,
	AuthPublicKeyFile: publicKeyFile,
	RegisterFunc: func(s *grpc.Server) {
		// 注意，在创建 GRPC 服务时，需要确定该服务是否涉及到“加密”：
		// 如有涉及，请传入 “TokenGenerator”。
		// “TokenGenerator”： 用于生成 JWT Token。

		// 注册etranslate服务
		etranslatepb.RegisterEtranslateServiceServer(s, &etranslate.Service{})

		// 注册auth服务
		authpb.RegisterAuthServiceServer(s, &auth.Service{
			TokenExpire:       appConfig.JwtInfo.Expire * time.Second, // token超时时间
			AuthPublicKeyFile: etRoot.Path("config/cert/public.key"),
			TokenGenerator:    token.NewJWTTokenGen(appConfig.JwtInfo.Issuer, privKey),
		})
	},
}); err != nil {
	zap.S().Panicf("cannot GRPC Run err: %v", err)
}

...

```

## 相关说明

1. 采用JWT授权登录，可参考```Auth:Login()```接口的生成token方法，前端获取到```token```后需要保存到```Header```中的```authorization```。
2. 登录状态判断，可参考```Auth:GetUserToken()```方法。
	- 用户请求接口时，gRPC拦截器会获取 Header 携带的```authorization```参数，解析后将用户```id```保存到上下下文中，通过上下文获取用户ID方法：```interceptor.UserIDFromContext(ctx)```
	```
	uid, err := auth.UserIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "用户未授权")
	}
	```
