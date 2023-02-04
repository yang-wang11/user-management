
#go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest   #v1.20

# generate go codes
function GenerateProtoCode() {
  OUT=$1
  protoc --go_opt=paths=source_relative --go_out=${OUT} --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative --go-grpc_out=${OUT} *.proto
  protoc-go-inject-tag -input="${OUT}/*.pb.go"
}

GenerateProtoCode ../api-gateway/pkg/pb
GenerateProtoCode ../user/pkg/pb
