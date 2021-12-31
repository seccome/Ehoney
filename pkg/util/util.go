package util

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/time_parse"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ip2location/ip2location-go/v9"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func GetCurrentTime() string {
	return time_parse.CSTLayoutString()
}

func ExecPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return path
}

func CheckInjectionData(payload string) bool {
	complite, _ := regexp.Compile(`^[a-zA-Z0-9\.\-\_\:\/\\]*$`)
	return !complite.MatchString(payload)
}

func GetLocationByIP(ip string) (*ip2location.IP2Locationrecord, error) {
	db, err := ip2location.OpenDB(path.Join(WorkingPath(), "data", "IP2LOCATION-LITE-DB.BIN"))
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	results, err := db.Get_all(ip)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &results, nil
}

func WorkingPath() string {
	path := os.Getenv("WorkingDir")
	if len(path) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return cwd
		}
		path = cwd
	}
	return path
}

func CheckPathIsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func RemoveDir(path string) error {
	if _, err := os.Stat(path); os.IsExist(err) {
		return os.RemoveAll(path)
	}
	return nil
}

func Base64Encode(str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return string(encoded)
}

func Base64Decode(str string) string {
	decoded, _ := base64.StdEncoding.DecodeString(str)
	return string(decoded)
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func GetCurrentFormatStr(fmtStr string) string {
	if fmtStr == "" {
		fmtStr = "2006-01-02 15:04:05"
	}
	return time.Now().Format(fmtStr)
}

func Sec2TimeStr(sec int64, fmtStr string) string {
	if fmtStr == "" {
		fmtStr = "2006-01-02 15:04:05"
	}
	return time.Unix(sec, 0).Format(fmtStr)
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func GetUniqueID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func createFlatZip(w io.Writer, files []string) error {
	z := zip.NewWriter(w)
	for _, file := range files {
		src, err := os.Open(file)
		if err != nil {
			return err
		}
		info, err := src.Stat()
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.Base(file)
		dst, err := z.CreateHeader(hdr)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		src.Close()
	}
	return z.Close()
}

func CompressZIP(destPath string, srcPath ...string) error {
	a, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer a.Close()

	err = createFlatZip(a, srcPath)
	if err != nil {
		return err
	}
	return nil
}

func createFlatTarGz(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		header := new(tar.Header)
		header.Name = filepath.Base(path)
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func CompressTarGz(destPath string, srcPath ...string) error {
	file, err := os.Create(destPath)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer file.Close()
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, path := range srcPath {
		if err := createFlatTarGz(tw, path); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetFileMD5(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func CopyFile(dstName, srcName string) (writeen int64, err error) {
	src, err := os.Open(dstName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(srcName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func DoEXEToken(sourceFile string, destFile string, traceUrl string) error {

	//TODO  local path and docker path is hard code FIXME
	dockerSourceFile := path.Join("/mnt/infile", path.Base(sourceFile))
	dockerDestFile := path.Join("/mnt/outfile", path.Base(sourceFile))
	localSourceFile := path.Join("/home/filetrace/infile", path.Base(sourceFile))
	localDestFile := path.Join("/home/filetrace/outfile", path.Base(sourceFile))
	fmt.Println("sourceFile", sourceFile)
	CreateDir(path.Dir(localSourceFile))
	CreateDir(path.Dir(localDestFile))
	CopyFile(sourceFile, localSourceFile)
	uri := "http://" + configs.GetSetting().Server.AppHost + ":5000/api/signmaking?urladdr=" + url.QueryEscape(traceUrl) + "&filetype=exe&inputfile=" + url.QueryEscape(dockerSourceFile) + "&outputfile=" + url.QueryEscape(dockerDestFile) + "&word="
	SendGETRequest(nil, uri)
	CreateDir(path.Dir(destFile))
	CopyFile(localDestFile, destFile)
	os.Remove(localSourceFile)
	os.Remove(localDestFile)
	return nil
}

type TokenFileCreateBody struct {
	SourceFile string
	DestFile   string
	TokenType  string
	TraceUrl   string

	Content         string
	ZipName         string // 设置压缩包的名字和路径
	ChildFolderName string // 压缩包里文件夹的名字
	ContainerPath   string // 压缩指定文件夹目录下的所有文件到压缩包里。(该文件夹下没有文件也可以,主要是为了可以添加一些其他诱饵文件)
	Host            string
	HostPort        string
	TraceCode       string
}

func CreateTokenFile(tokenFileCreateBody TokenFileCreateBody) error {
	zap.L().Info(fmt.Sprintf("开始注入文件[%s] 类型[%s]蜜签", tokenFileCreateBody.TokenType, tokenFileCreateBody.SourceFile))
	if tokenFileCreateBody.TokenType != "BrowserPDF" && tokenFileCreateBody.TokenType != "WPS" && !CheckPathIsExist(tokenFileCreateBody.SourceFile) {
		zap.L().Error("待加签文件不存在")
		return errors.New("source file is not exist")
	}
	var toolPath = path.Join(WorkingPath(), "tool", "file_token_trace", "linux", "TraceFile")

	if !CheckPathIsExist(toolPath) {
		zap.L().Error("加签工具不存在")
		return errors.New("trace file is not exist")
	}
	CreateDir(path.Dir(tokenFileCreateBody.DestFile))

	var cmd *exec.Cmd

	// 命令组装
	switch tokenFileCreateBody.TokenType {

	case "WPS":
		cmd = exec.Command(toolPath, "-u", tokenFileCreateBody.TraceUrl, "-o", tokenFileCreateBody.DestFile, "-w", tokenFileCreateBody.Content, "-t", "wps")
	case "BrowserPDF":
		cmd = exec.Command(toolPath, "-u", tokenFileCreateBody.TraceUrl, "-o", tokenFileCreateBody.DestFile, "-w", tokenFileCreateBody.Content, "-t", "chromepdf")
	case "FILE":
		cmd = exec.Command(toolPath, "-u", tokenFileCreateBody.TraceUrl, "-o", tokenFileCreateBody.DestFile, "-i", tokenFileCreateBody.SourceFile, "-t", "office")
	case "EXE":
		cmd = exec.Command(toolPath, "-u", tokenFileCreateBody.TraceUrl, "-o", tokenFileCreateBody.DestFile, "-i", tokenFileCreateBody.SourceFile, "-t", "exe")
	//case "WIN_FOLDER":
	//	cmd = exec.Command(toolPath, "--zn", tokenFileCreateBody.zipName, "--fn", tokenFileCreateBody.childFolderName, "--dp", tokenFileCreateBody.containerPath, "--hn", tokenFileCreateBody.host, "--hp", tokenFileCreateBody.hostPort,
	//		"--tc", tokenFileCreateBody.traceCode, "-t", "winfolder")
	default:
		zap.L().Error("无法处理的蜜签类型: " + tokenFileCreateBody.TokenType)
		return errors.New("无法处理的蜜签类型: " + tokenFileCreateBody.TokenType)
	}
	cmd.Dir = path.Dir(toolPath)
	zap.L().Info(tokenFileCreateBody.TraceUrl)
	zap.L().Info(tokenFileCreateBody.DestFile)
	zap.L().Info(tokenFileCreateBody.SourceFile)
	zap.L().Info(tokenFileCreateBody.Content)

	zap.L().Info("cmd : " + cmd.String())

	_, err := cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("文件密签加签失败")
		zap.L().Error(err.Error())
		fmt.Println("文件密签加签失败:" + err.Error())
		os.RemoveAll(path.Dir(tokenFileCreateBody.DestFile))
		return err
	}
	return nil
}

func DoFileTokenTrace(sourceFile string, destFile string, traceUrl string) error {
	if !CheckPathIsExist(sourceFile) {
		zap.L().Error("待加签文件不存在")
		return errors.New("source file is not exist")
	}
	var toolPath = path.Join(WorkingPath(), "tool", "file_token_trace", runtime.GOOS, "TraceFile")
	if !CheckPathIsExist(toolPath) {
		zap.L().Error("加签工具不存在")
		return errors.New("trace file is not exist")
	}
	CreateDir(path.Dir(destFile))

	//TODO support mac version
	cmd := exec.Command(toolPath, "-i", sourceFile, "-o", destFile, "-u", traceUrl, "-t", traceUrl)
	cmd.Dir = path.Dir(toolPath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("文件密签加签失败")
		zap.L().Error(err.Error())
		fmt.Println("文件密签加签失败:" + err.Error())
		os.RemoveAll(path.Dir(destFile))
		return err
	}
	return nil
}

func DoBrowserPDFToken(data string, destFile string, traceUrl string) error {
	var toolPath = path.Join(WorkingPath(), "tool", "file_token_trace", runtime.GOOS, "TraceBrowserPDF")
	if !CheckPathIsExist(toolPath) {
		zap.L().Error("BrowserPDF密签工具不存在")
		return errors.New("BrowserPDF密签工具不存在")
	}
	CreateDir(path.Dir(destFile))
	cmd := exec.Command(toolPath, "-w", data, "-o", destFile, "-u", traceUrl)
	cmd.Dir = path.Dir(toolPath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("BrowserPDF密签生成失败")
		zap.L().Error(err.Error())
		fmt.Println("BrowserPDF密签生成失败:" + err.Error())
		os.RemoveAll(path.Dir(destFile))
		return err
	}
	return nil
}

func SendGETRequest(header map[string]string, uri string) ([]byte, error) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", uri, nil)

	//add header
	for key, value := range header {
		request.Header.Add(key, value)
	}
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	ret, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ret, nil
}

func TcpGather(ips []string, port string) bool {
	for _, ip := range ips {
		fmt.Println("test ip:" + ip)
		address := net.JoinHostPort(ip, port)
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			return false
		} else {
			if conn != nil {
				return true
				_ = conn.Close()
			} else {
				return false
			}
		}
	}
	return false
}

func SendDingMsg(title, name, msg string) error {
	webHook := configs.GetSetting().App.WebHook
	content := `{"msgtype": "markdown", "markdown":{"title":"` + title + `","text": "### ` + name + `\n > ` + msg + `"}}`
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		zap.L().Error("发送请求失败")
		return err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client.Do(req)
	return nil
}

func IsLocalIP(ip string) bool {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for i := range address {
		intf, _, err := net.ParseCIDR(address[i].String())
		if err != nil {
			return false
		}
		if net.ParseIP(ip).Equal(intf) {
			return true
		}
	}
	return false
}
