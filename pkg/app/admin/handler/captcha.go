package handler

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/log"
)

func GetCaptcha(c *gin.Context) {
	captchaId := captcha.New()
	c.SetCookie(CookieKeyCaptchaId, captchaId, 300, "", "", http.SameSiteDefaultMode, false, false)

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")

	var content bytes.Buffer
	c.Writer.Header().Set("Content-Type", "image/png")
	err := captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return
	}
	http.ServeContent(c.Writer, c.Request, captchaId+"png", time.Time{}, bytes.NewReader(content.Bytes()))
}

// curl http://localhost:8080/verifyCaptcha?captchaId=OzXt7zGEeq27xtGvD2fC&value=507826
func VerifyCaptcha(c *gin.Context, value string) bool {
	requestId := c.GetString(CtxKeyRequestId)
	captchaId, err := c.Cookie(CookieKeyCaptchaId)
	if err != nil {
		log.Trace(requestId).Errorf("cookie里没有_captchaId")
		return false
	}

	if captchaId == "" || value == "" {
		log.Errorf("没有找到验证码captchaId(%s)或者验证码value(%s)", captchaId, value)
		return false
	}

	if captcha.VerifyString(captchaId, value) {
		log.Trace(requestId).Infof("验证码验证成功,captchaId:%s value:%s", captchaId, value)
		return true
	}

	log.Trace(requestId).Errorf("验证码验证失败,captchaId:%s value:%s", captchaId, value)
	return false
}
