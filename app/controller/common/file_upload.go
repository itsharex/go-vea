package common

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/global"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileApi struct{}

func (*FileApi) UploadFile(c *gin.Context) {
	// 限制上传的文件大小，最大为 8MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8*1024*1024)

	// 解析 form 表单，获取上传的文件
	if err := c.Request.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}
	defer file.Close()

	// 将上传的文件保存到本地
	datePath := time.Now().Format("20060102") // 格式为yyyyMMdd
	filePath := "uploads/" + datePath + "/" + header.Filename
	dirPath := filepath.Dir(filePath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		// 处理创建目录时发生的错误
		global.Logger.Error(err)
	}
	dst, err := os.Create(filePath)
	if err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	result.Ok(c)
}
