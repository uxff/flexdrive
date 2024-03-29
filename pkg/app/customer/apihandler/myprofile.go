package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/log"
)

func init() {
}

func Profile(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	loginInfo := getLoginInfo(c)
	levelInfo, err := dao.GetUserLevelById(loginInfo.UserEnt.LevelId)
	if err != nil {
		log.Trace(requestId).Errorf("query level error:%v", err)
		JsonErrMsg(c, ErrInvalidParam, "查询等级错误")
		return
	}

	if levelInfo == nil {
		levelInfo, err = dao.GetDefaultUserLevel()
		if err != nil {
			log.Trace(requestId).Errorf("query level error:%v", err)
			JsonErrMsg(c, ErrInvalidParam, "查询等级错误")
			return
		}
	}

	if levelInfo == nil {
		JsonErrMsg(c, ErrInvalidParam, "无默认等级")
		return
	}

	JsonOk(c, levelInfo)
}
