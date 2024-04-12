# build
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/tools/drone_plugin -e GOPROXY=$GOPROXY $GOIMAGE go build -o /work/tools/drone_plugin/dingding_notify/notify /work/tools/drone_plugin/dingding_notify
docker build -t jybl/notify tools/drone_plugin/dingding_notify
docker push jybl/notify


# test notify
docker run --rm -e PLUGIN_DING_TOKEN=xxx -e PLUGIN_DING_SECRET=xxx jybl/notify
