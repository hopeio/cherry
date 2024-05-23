gopath=/mnt/d/SDK/gopath
if [ -z "$1" ]; then
  echo "gopath参数为空"
  exit 1
else
  gopath=$1
  echo "gopath: $gopath"
fi
protoc=/mnt/d/tools/protoc-22.3-linux-x86_64
if [ -z "$2" ]; then
  echo "protoc参数为空"
  exit 1
else
  protoc=$2
  echo "protoc: $protoc"
fi
cd $(dirname $0) && pwd
cherry_dir=$(cd ../../..;pwd)

goproxy=https://goproxy.io,https://goproxy.cn,direct
goimage=golang:latest

if [ -n "$3" ]; then
  goimage= $1
fi

# install tools
docker run --rm -e GOPROXY=$goproxy -e GOFLAGS=-buildvcs=false -v $gopath:/go -v $protoc:/protoc -v $cherry_dir:/work -w /work/tools/protoc --name install $goimage bash ./install_tools.sh /protoc
# docker rm -f install
echo "docker build"
dockerTmpDir=$cherry_dir/tools/protoc/_docker
mkdir $dockerTmpDir
cp $gopath/bin/protoc-gen-enum $dockerTmpDir/
cp $gopath/bin/protoc-gen-go $dockerTmpDir/
cp $gopath/bin/protoc-gen-go-grpc $dockerTmpDir/
cp $gopath/bin/protoc-gen-go-patch $dockerTmpDir/
cp $gopath/bin/protoc-gen-validator $dockerTmpDir/
cp $gopath/bin/protoc-gen-grpc-gateway $dockerTmpDir/
cp $gopath/bin/protoc-gen-grpc-gin $dockerTmpDir/
cp $gopath/bin/protoc-gen-openapiv2 $dockerTmpDir/
cp $gopath/bin/protoc-gen-gql $dockerTmpDir/
cp $gopath/bin/protoc-gen-gogql $dockerTmpDir/
cp $gopath/bin/gqlgen $dockerTmpDir/
cp $gopath/bin/protogen $dockerTmpDir/
cp -r $cherry_dir/protobuf/_proto $dockerTmpDir/_proto
cp -r $protoc $dockerTmpDir/protoc


docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $cherry_dir/tools/internal/docker/Dockerfile_local $dockerTmpDir
rm -rf $dockerTmpDir
docker push jybl/protogen