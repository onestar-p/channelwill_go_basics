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

1. 在```channelwill_go_basics/proto```目录下创建```.proto```文件
2. 在```channelwill_go_basics/proto```目录下创建```.yaml```文件，用于规范```RESTful api```.


例子：

auth.proto
```
syntax = "proto3";
import "google/protobuf/empty.proto";
package auth.v1;
option go_package="channelwill_go_basics/proto/auth/v1;authpb";

service AuthService {
    rpc Login (google.protobuf.Empty) returns (LoginResponse);
    rpc GetUserToken (google.protobuf.Empty) returns (GetUserTokenResponse);
}

message LoginResponse {
    string user_token = 1;
}

message GetUserTokenResponse {
    string token  = 1;
}
```

auth.yaml
```
type: google.api.Service
config_version: 3

http: 
  rules:
    - selector: auth.v1.AuthService.Login
      post: /v1/login
      body: "*"
    - selector: auth.v1.AuthService.GetUserToken
      post: /v1/get_user_token
      body: "*"
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
- 文件：```channelwill_go_basics/utils/register/server/server.go```
- 方法：```func NewClientServices() []*ClientService```
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
1. 创建接口文件
	> 在目录```channelwill_go_basics/service```下创建对应接口目录
2. 注册服务
	> 在```root/cmd/server/main.go```文件注册创建好的接口

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
				TokenGenerator:    jwt.NewJWTTokenGen(appConfig.JwtInfo.Issuer, privKey),
			})
		},
	}); err != nil {
		zap.S().Panicf("cannot GRPC Run err: %v", err)
	}

	...

	```

### 四、其他说明
1. 生成JWT Token
```
tokenExpire := 7200 * time.Second // 失效时间
tkn, err := s.TokenGenerator.GenerateToken(aid, tokenExpire)
if err != nil {
	zap.S().Error("cannot generate token", zap.Error(err))
	return nil, status.Error(codes.Internal, "")
}
fmt.Println(tkn)
```
2. 登录状态判断
> 前端请求需要把```token```信息存放到```Header authorization```。\
> 后端通过上下文保存用户ID
```
c := context.Background()
uid, err := interceptor.UserIDFromContext(c)
if err != nil {
	return nil, status.Error(codes.Unauthenticated, "用户未授权")
}
fmt.Println(uid)
```


## 本地包说明
1. channelwill_go_basics/utils/jwt：
	- ```jwt.NewJWTKey(privateKeyFilePath).GetPrivateKey()```：获取私钥
	- ```jwt.NewJWTKey(publicKeyFilePath).GetPublicKey()```：获取公钥
	- ```jwt.NewJWTTokenGen("test", privKey).GenerateToken("id", 7200*time.Second)```：生成token
	- ```jwt.NewJWTTokenVerifyer(pubKey).Verify(tkn)```：验证秘钥，得到解析结果

2. channelwill_go_basics/utils/auth:
	- ```auth.NewAuthContext().ContextWithUserId(c, uid)```：用户ID写入上下文
	- ```auth.NewAuthContext().UserIDFromContext(c)```：通过上下文获取用户ID
