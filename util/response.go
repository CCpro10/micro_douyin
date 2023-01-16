package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponse struct {
	StatusCode int // 状态码，0-成功，其他值-失败
	StatusMsg  string
	ReturnVal  map[string]interface{}
}

func MakeResponse(c *gin.Context, resp *HttpResponse) {
	var msg string
	if resp.StatusMsg != "" {
		msg = resp.StatusMsg
	} else {
		msg = ErrCode2Msg[resp.StatusCode]
	}

	respMap := gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  msg,
	}
	if len(resp.ReturnVal) != 0 {
		for k, v := range resp.ReturnVal {
			respMap[k] = v
		}
	}

	c.JSON(http.StatusOK, respMap)
}
