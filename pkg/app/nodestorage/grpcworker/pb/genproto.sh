# 将proto生成go代码 
# 请前提安装好2个工具: protoc 和 protoc-gen-go
protoc -I pingablepb/ pingablepb/pingable.proto --go_out=plugins=grpc:pingablepb