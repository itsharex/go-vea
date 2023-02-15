package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-vea/app/common/constant"
	"go-vea/app/core"
	"go-vea/app/model/monitor"
	"go-vea/app/service/monitorsrv"
	"go-vea/global"
	"go-vea/util"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func OperationRecord(title string, businessType constant.BusinessType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		// 请求body
		var reqBody []byte
		if ctx.Request.Method != http.MethodGet {
			reqBody, err = io.ReadAll(ctx.Request.Body)
			if err != nil {
				global.Logger.Error("read body from request error:", err)
			} else {
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}
		} else {
			query := ctx.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			reqBody, _ = json.Marshal(&m)
		}

		// 响应body
		writer := responseBodyWriter{
			ResponseWriter: ctx.Writer,
			body:           &bytes.Buffer{},
		}
		ctx.Writer = writer

		handler := ctx.HandlerName()
		loginUser, _ := core.TokenSrv.GetLoginUser(ctx)
		now := time.Now()
		sysOperLog := monitor.SysOperLog{
			OperIP:        ctx.ClientIP(),
			Title:         title,
			OperName:      loginUser.SysUserResp.SysUser.Username,
			DeptName:      loginUser.SysUserResp.SysDept.DeptName,
			BusinessType:  int64(businessType),
			OperatorType:  constant.OPER_MANAGE,
			RequestMethod: ctx.Request.Method,
			OperURL:       ctx.Request.URL.Path,
			Status:        util.StatusConvert(ctx.Writer.Status()),
			Method:        handler,
			OperTime:      &now,
			OperParam:     string(reqBody),
			JSONResult:    writer.body.String(),
		}
		if err != nil {
			sysOperLog.ErrorMsg = err.Error()
		}
		// Next() 进入下一个middleware
		ctx.Next()
		latency := time.Since(now)
		// IP属地
		region, err := global.IpSearcher.SearchByStr(ctx.ClientIP())
		sysOperLog.OperLocation = region
		fmt.Println(latency)

		err = monitorsrv.SysOperLogSrv.AddSysOperLog(ctx, &sysOperLog)
		if err != nil {
			global.Logger.Error("oper log err")
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
