// 更新用户本地配置文件
// 1. 生成指定目录下的每个文件的hash值，并保存为json文件，这个json文件也会上传到服务器
// 2. 定时
//		 a. 扫描配置目录按照目录结构生成各个文件的hash，并保存为json文件
//       b. 比较服务器上的json文件的hash与本地json文件的hash，如果一致，不用更新
//		 c. 如果不一致, 下载服务器的json文件与本地的json文件对比，有hash不一样的文件需要被上传至服务器进行更新

// 需要得到eagle_wing的cfg目录的位置
// windows获得文件hash的方法 Get-FileHash -path "D:\github\NATS\fork\go-nats\examples\cloud_ftp\2.json" -Algorithm sha1
package main
import (
	"flag"
	"log"
	"os"
	"io"
	"fmt"
	"path/filepath"
	"crypto/sha1"
	"encoding/json"
)

var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Debug *log.Logger
var cfgPath = flag.String("f", "../cfg", "http service address")
type HashFileStruct struct {
	FileName  string
	Sha1Value string
}
var FileHash = make([]HashFileStruct,0)
func init() {
	file , err := os.OpenFile("cloud-service.log",os.O_APPEND |os.O_CREATE|os.O_WRONLY,0666)
	//如果打开错误日志文件失败
	if err != nil{
		log.Fatalln("打开日志文件失败。。")
	}
	//初始化错误日志记录器
	Info = log.New(io.MultiWriter(os.Stderr,file),"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	Warn = log.New(io.MultiWriter(os.Stderr,file),"Warn:",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr,file),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
	Debug = log.New(file,"Debug:",log.Ldate | log.Ltime | log.Lshortfile)
}

func getCfgDir() (string, error){
	flag.Parse()
	return filepath.Abs(*cfgPath)
}

func walkHandler( path string, f os.FileInfo, err error ) error {
	sha1Value, _ := getSHA1OfFile(path)
	FileHash = append(FileHash, HashFileStruct{path, sha1Value})
	return nil
}

func getSHA1OfFile(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "",err
	}

	h := sha1.New()
	_, err = io.Copy(h,file)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%x",h.Sum(nil)), nil
}

func main(){
	cfgPath, err := getCfgDir()
	if err!=nil{
		log.Fatalf("Get config dir fail, %s", err.Error())
	}
	fmt.Printf("%s\n", cfgPath)
	filepath.Walk(cfgPath, walkHandler)
	j, err := json.Marshal(FileHash)
	fmt.Printf("%v", string(j))
	//for _,item := range FileHash{
	//	fmt.Printf("%s = %s\n", item.FileName, item.Sha1Value)
	//}
}