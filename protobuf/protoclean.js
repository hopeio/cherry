import fs from "fs";
import path from "path";

function main() {
    let inputpath = '../../../proto'
    let outputpath = '../../../proto_std'
    clean(inputpath,outputpath)
}

function clean(inputpath, outputpath) {
    console.log(inputpath, "[", outputpath, "]")
    if (!fs.existsSync(outputpath)) fs.mkdirSync(outputpath)
    const files = fs.readdirSync(inputpath)
    for (let file of files) {
        const newPath = path.join(inputpath, file)
        fs.stat(newPath, function (err, stats) {
            if (stats.isFile()) replace(newPath, path.join(outputpath, file))
            else if (newPath.endsWith(`utils/proto`)) return
            else clean(newPath, path.join(outputpath, file))
        })
    }
}

function replace(inputfile, outputfile) {
    if (!inputfile.endsWith(".proto")) return
/*    const buf = Buffer.alloc(1024);
    fs.open(inputfile, "r+", function (err, fd) {
        if (err) return console.error(err);
        console.log("文件打开成功！");
        console.log("准备读取文件！");
        fs.read(fd, buf, 0, buf.length, 0, function (err, bytes) {
            if (err) console.log(err);
            console.log(bytes + "  字节被读取");

            // 仅输出读取的字节
            if (bytes > 0) {
                console.log(buf.slice(0, bytes).toString());
            }
        });
        // 关闭文件
        fs.close(fd, function (err) {
            if (err) console.log(err);
            console.log("文件关闭成功");
        });
    })*/
    let data = fs.readFileSync(inputfile,'utf-8')
    data = data.replace(/import "github.*\n/g,'')
    data = data.replace(/import "protoc-gen-openapiv2.*\n/g,'')
    data = data.replace(/import "utils\/proto\/.*\n/g,'')
    data = data.replace(/import "google\/api\/annotations.proto.*\n/g,'')
    data = data.replace(/\[\([\w.]*\)[\\\n\w=/."',:;*@(){} \-\u4e00-\u9fa5]*]/g,'')
    data = data.replace(/option \([\w.]*\)[\[\]\\\n\w=/."',:*@(){} \-\u4e00-\u9fa5]*;/g,'')
    fs.writeFileSync(outputfile,data)
}

main()
