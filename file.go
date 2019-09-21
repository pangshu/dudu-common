package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// 判断是否为目录
func (*DuduFile) IsDir(path string) bool {
	tmpInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	if nil != err {
		return false
	}

	return tmpInfo.IsDir()
}

// 根据路据创建文件夹
func (*DuduFile) Mkdir(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return false
		}
		// 检查文件夹是否创建成功
		if _, err := os.Stat(path); err == nil {
			return true
		} else {
			return false
		}
	}
}


// 获取文件类型
func (*DuduFile) GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// 获取文件大小
func (*DuduFile) GetFileSize(path string) (int64,error) {
	fi, err := os.Stat(path)
	if nil != err {
		return -1,err
	}

	return fi.Size(),nil
}

// 判断文件是否存在
func (*DuduFile) IsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}

// 判断文件是否为二进制文件
func (*DuduFile) IsBinary(content string) bool {
	for _, b := range content {
		if 0 == b {
			return true
		}
	}

	return false
}

// 复制文件到目标地址
func (*DuduFile) CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil && sourceinfo != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return nil
}

// 复制文件夹到指定位置
func (*DuduFile) CopyDir(source string, dest string) (err error) {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	defer directory.Close()

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		srcFilePath := filepath.Join(source, obj.Name())
		destFilePath := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// create sub-directories - recursively
			err = File.CopyDir(srcFilePath, destFilePath)
			if err != nil {
				return err
			}
		} else {
			err = File.CopyFile(srcFilePath, destFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


// 获取项目路径
func (*DuduFile) GetProjectPath() string{
	var projectPath string
	projectPath, _ = os.Getwd()
	return projectPath
}

// 获取config路径
func (*DuduFile) GetConfigPath() string{
	path,_ := os.Getwd()
	var osType = runtime.GOOS
	if osType == "windows"{
		path = path + "\\" + "config\\"
	}else if osType == "linux"{
		path = path +"/" + "config/"
	}
	return  path
}

//获取应用文件名
func (*DuduFile) GetAppName() string {
	full := os.Args[0]
	full = strings.Replace(full, "\\", "/", -1)
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		name = strings.TrimSuffix(name, ".exe")
		return name
	}

	return ""
}

//获取当前路径目录
func (*DuduFile) GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Print(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
