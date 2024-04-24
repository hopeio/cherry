cd $(dirname $0) && pwd
cherry_dir=$(cd ../../..;pwd)

goimage=golang:latest

if [ -n "$1" ]; then
  goimage= $1
fi

docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $cherry_dir/tools/protoc/internal/Dockerfile_upgrade .