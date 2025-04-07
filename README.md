# cherry
cherry是一个基于go语言的http&grpc框架，可通过一个实现一个端口同时提供http和grpc服务

## grpc-gateway
protobuf插件生成基于gin的gateway,实现基于grpc的实现对外提供http的接口
## grpc-web
grpc-web集成,通过浏览器调用grpc服务
## http3
通过quic实现http3的支持
## 可观测性
集成opentelemetry
## ~~graphql~~(废弃)
鉴于国内实际落地并不多，移除掉graphql支持