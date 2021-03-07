package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/utils"
	"os"
	"path"
	"strings"
	"time"
)

func (svr *Server) handleUploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		Err(c, err.Error())
		return
	}

	extname := strings.ToLower(path.Ext(file.Filename))

	if !utils.IsImageExt(extname) {
		Err(c, "UnSupport upload file type "+extname)
		return
	}
	maxUploadSize := 10485760 //10M
	if file.Size > int64(maxUploadSize) {
		Err(c, "Upload file too large")
		return
	}
	uid := uuid.NewV4()
	name := uid.String() + path.Ext(file.Filename)
	now := time.Now()
	path := fmt.Sprintf("./static/file/%d/%d/%d", now.Year(), now.Month(), now.Day())
	err = os.MkdirAll(path, 666)
	if err != nil {
		Err(c, "store file error")
		logrus.Errorf("MkdirAll err %s", err)
		return
	}
	var dst = fmt.Sprintf("%s/%s", path, name)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		Err(c, err.Error())
		logrus.Errorf("SaveUploadedFile err %s", err.Error())
		return
	}
	Ok(c, fmt.Sprintf("/static/file/%d/%d/%d/%s", now.Year(), now.Month(), now.Day(), name))
}
