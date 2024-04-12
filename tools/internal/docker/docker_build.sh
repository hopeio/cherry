cd $(dirname $0) && pwd
cherry_dir=$(cd ../../..;pwd)
dockerTmpDir=$cherry_dir/tools/protoc/_docker
cd $dockerTmpDir
docker build -t jybl/goprotoc --build-arg IMAGE=$goimage -f $cherry_dir/tools/internal/Dockerfile_local $dockerTmpDir
rm -r $dockerTmpDir
docker push jybl/goprotoc