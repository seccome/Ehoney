package util

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func GetCurrentPathString() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path = string(path[0:strings.LastIndex(path, "/")])
	return path
}

func ModifyFile(newpath string, newbaitpath string, baitname string, baitfilename string) {
	pocpath := GetCurrentPathString() + "/policy/"
	fileName := "install.sh"
	pocpath = pocpath + fileName
	in, err := os.Open(pocpath)
	if err != nil {
		logs.Error("ModifyFile open file fail:", err)
		fmt.Println("open file fail:", err)
		//os.Exit(-1)
	}
	defer in.Close()

	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return
	}
	out, err := os.OpenFile(newpath+"/"+fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		logs.Error("ModifyFile Open write file fail:", err)
		fmt.Println("Open write file fail:", err)
		//os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	// index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logs.Error("ModifyFile read err:", err)
			fmt.Println("read err:", err)
			//os.Exit(-1)
		}
		if find := strings.Contains(string(line), "7779f3be0bbddcf3d9f4870d44629681"); find {
			newLine := strings.Replace(string(line), "7779f3be0bbddcf3d9f4870d44629681", newbaitpath+"/"+baitfilename, -1)
			_, err = out.WriteString(newLine + "\n")
			if err != nil {
				logs.Error("ModifyFile write to file fail:", err)
				fmt.Println("write to file fail:", err)
			}
		} else if find := strings.Contains(string(line), "6669f3be0bbddcf3d9f4870d44629681"); find {
			newLine := strings.Replace(string(line), "6669f3be0bbddcf3d9f4870d44629681", baitfilename, -1)
			_, err = out.WriteString(newLine + "\n")
			if err != nil {
				logs.Error("ModifyFile write to file fail:", err)
				fmt.Println("write to file fail:", err)
			}
		} else {
			_, err = out.WriteString(string(line) + "\n")
			if err != nil {
				logs.Error("ModifyFile write to file fail:", err)
				fmt.Println("write to file fail:", err)
			}
		}
		// fmt.Println("done ", index)
		// index++
	}
	fmt.Println("FINISH!")

}

func ModifySignUninstallFile(newpath string, baitname string) {
	pocpath := GetCurrentPathString() + "/policy/uninstall/"
	fileName := "install.sh"
	pocpath = pocpath + fileName
	in, err := os.Open(pocpath)
	if err != nil {
		logs.Error("open file fail:", err)
		//os.Exit(-1)
	}
	defer in.Close()
	err = os.MkdirAll("upload/honeysign/"+baitname, os.ModePerm)
	if err != nil {
		return
	}

	out, err := os.OpenFile("upload/honeysign/"+baitname+"/"+fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		logs.Error("Open write file fail:", err)
		//os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	// index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logs.Error("read err:", err)
			//os.Exit(-1)
		}
		newLine := strings.Replace(string(line), "8889f3be0bbddcf3d9f4870d44629681", newpath, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			logs.Error("write to file fail:", err)
			//os.Exit(-1)
		}
		// fmt.Println("done ", index)
		// index++
	}
	fmt.Println("FINISH!")

}

func ModifyUninstallFile(newpath string, baitname string) {
	pocpath := GetCurrentPathString() + "/policy/uninstall/"
	fileName := "install.sh"
	pocpath = pocpath + fileName
	in, err := os.Open(pocpath)
	if err != nil {
		logs.Error("[ModifyUninstallFile] open file fail:", err)
	}
	defer in.Close()
	err = os.MkdirAll("upload/"+baitname, os.ModePerm)
	if err != nil {
		return
	}

	out, err := os.OpenFile("upload/"+baitname+"/"+fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		logs.Error("[ModifyUninstallFile] Open write file fail:", err)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logs.Error("[ModifyUninstallFile] read err:", err)
		}
		newLine := strings.Replace(string(line), "8889f3be0bbddcf3d9f4870d44629681", newpath, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			logs.Error("[ModifyUninstallFile] write to file fail:", err)
		}
	}
	fmt.Println("FINISH!")

}

func GetPercent(a int, b int) string {
	num1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(a)/float64(b)), 64)
	pernum := strconv.FormatFloat(num1*100, 'f', 0, 64)
	return pernum + "%"
}

func Gzip(filepath, filename string) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()
	gw := gzip.NewWriter(File)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return walk(filepath, tw)
}

func walk(path string, tw *tar.Writer) error {
	path = strings.Replace(path, "\\", "/", -1)
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	index := strings.Index(path, "/")
	list := strings.Join(strings.Split(path, "/")[index:], "/")
	for _, v := range info {
		if v.IsDir() {
			head := tar.Header{Name: list + v.Name(), Typeflag: tar.TypeDir, ModTime: v.ModTime()}
			tw.WriteHeader(&head)
			walk(path+v.Name(), tw)
			continue
		}
		F, err := os.Open(path + v.Name())
		if err != nil {
			fmt.Println("打开文件%s失败.", err)
			continue
		}
		head := tar.Header{Name: list + v.Name(), Size: v.Size(), Mode: int64(v.Mode()), ModTime: v.ModTime()}
		tw.WriteHeader(&head)
		io.Copy(tw, F)
		F.Close()
	}
	return nil
}

func FileTarZip(sourcepath string, destpath string) {
	fw, err := os.Create(destpath)
	if err != nil {
		logs.Error("[FileTarZip]err:", err)
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打开文件夹
	dir, err := os.Open(sourcepath)
	if err != nil {
		logs.Error("[FileTarZip]err:", err)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		logs.Error("[FileTarZip]err:", err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			continue
		}

		// 打印文件名称
		fmt.Println(fi.Name())

		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			logs.Error("[FileTarZip]err:", err)
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			logs.Error("[FileTarZip]err:", err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			logs.Error("[FileTarZip]err:", err)
		}
	}

	fmt.Println("tar.gz ok")
}

func MakeDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func Base64Encode(str string) string {
	strbytes := []byte(str)
	encoded := base64.StdEncoding.EncodeToString(strbytes)
	return string(encoded)
}

func Base64Decode(str string) string {
	//strbytes := []byte(str)
	decoded, _ := base64.StdEncoding.DecodeString(str)
	return string(decoded)
}

func GetFileMd5(filepath string) string {
	file, _ := os.Open(filepath)
	md5 := md5.New()
	io.Copy(md5, file)
	MD5Str := hex.EncodeToString(md5.Sum(nil))
	return MD5Str
}

func GetStrMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

func UntarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
		//clamavdockerpath := beego.AppConfig.String("clamavdockerpath")
		//clamavdockername := beego.AppConfig.String("clamavdockername")
		//if clamavdockername != ""{
		//	cmd := exec.Command("docker", "cp", destFileName, clamavdockerpath+":"+clamavdockername+"/"+ header.Name[len(prefix):])
		//	out, err := cmd.CombinedOutput()
		//	if err != nil {
		//		logs.Error("docker cp err:",err)
		//	}
		//	if len(out) != 0 {
		//		log.Print("docker cp :",string(out))
		//	}
		//}
	}

	return nil
}

func GetIP(srcip string) string {
	resultip := ""
	ipRegexp := regexp.MustCompile(`^?([^:]*)`)
	ipparams := ipRegexp.FindStringSubmatch(srcip)
	if len(ipparams) > 0 {
		resultip = ipparams[0]
	}
	return resultip
}

func GetHost(srcip string) string {
	resultip := ""
	ipRegexp := regexp.MustCompile(`^(http://)(.*)(/)$`)
	params := ipRegexp.FindStringSubmatch(srcip)
	if params != nil && len(params) > 3 {
		resultip = params[2]
	}
	return resultip
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func FindInList(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func NetWorkStatus(ip string) bool {
	cmd := exec.Command("ping", ip, "-c", "1", "-W", "5")
	err := cmd.Run()
	if err != nil {
		fmt.Println("NetWorkStatus Error:", ip, err.Error())
		return false
	}
	return true
}

func NetConnectTest(host string, port string) bool {
	timeout := 3 * time.Second
	iscon := false
	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		iscon = false
		//_, err_msg := err.Error()[0], err.Error()[5:]
	} else {
		iscon = true
		conn.Close()
	}
	return iscon
}


func DoStr(str string) string {
	var a = []byte(str)
	for i, ch := range a {

		switch {
		case ch > '~':
			a[i] = ' '
		case ch == '\r':
		case ch == '\n':
		case ch == '\t':
		case ch < ' ':
			a[i] = ' '
		}
	}
	return string(a)
}