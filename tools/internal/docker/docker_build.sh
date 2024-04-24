cd $(dirname $0) && pwd
cherry_dir=$(cd ../../..;pwd)

goimage=golang:latest
if [ -n "$1" ]; then
  goimage= $1
fi

dockerTmpDir=$cherry_dir/tools/protoc/_docker
mkdir $dockerTmpDir && cd $dockerTmpDir
docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $cherry_dir/tools/internal/docker/Dockerfile $dockerTmpDir
rm -rf $dockerTmpDir
docker push jybl/protogen