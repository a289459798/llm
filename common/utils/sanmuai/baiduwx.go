package sanmuai

import (
	"bytes"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

type BaiduWX struct {
	Ctx        context.Context
	SvcCtx     *svc.ServiceContext
	HTTPClient *http.Client
	Token      string
	AppKeyId   int
}

type ErrorResponse struct {
	Error *struct {
		Code    *int    `json:"code,omitempty"`
		Message string  `json:"message"`
		Param   *string `json:"param,omitempty"`
		Type    string  `json:"type"`
	} `json:"error,omitempty"`
}

type BDTokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type BDImageTaskResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data BDImageTask `json:"data"`
}

type BDImageTaskRequest struct {
	Prompt string                `json:"prompt"`
	Image  *multipart.FileHeader `json:"image"`
}

type BDImageTask struct {
	TaskId   int `json:"taskId"`
	AppKeyId int `json:"appKeyId"`
}

type BDImageTaskResultResponse struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data BDImageTaskResult `json:"data"`
}

type BDImageTaskResult struct {
	Status     int                 `json:"status"`
	ImgUrls    []map[string]string `json:"imgUrls"`
	CreateTime string              `json:"createTime"`
	Waiting    string              `json:"waiting"`
}

func NewBaiduWX(c context.Context, svcCtx *svc.ServiceContext) *BaiduWX {
	return &BaiduWX{
		Ctx:        c,
		SvcCtx:     svcCtx,
		Token:      "",
		HTTPClient: &http.Client{},
	}
}

func (b *BaiduWX) getClient(id int) error {
	apikey := &model.Apikey{}
	db := b.SvcCtx.Db.Where("channel = ?", "baidu")
	if id > 0 {
		db.Where("id = ?", id)
	} else {
		db.Where("status = ?", 1).Order("rand()")
	}
	db.Limit(1).Find(apikey)
	if apikey.Token == "" || time.Now().Unix()-apikey.UpdateAt.Unix() > 85400 {
		token, err := b.getToken(apikey.Key, apikey.Secret)
		if err != nil {
			return err
		}
		if token.Code != 0 {
			return errors.New(token.Msg)
		}
		b.Token = token.Data
		apikey.Token = token.Data
		apikey.UpdateAt = time.Now()
		b.SvcCtx.Db.Save(&apikey)
	} else {
		b.Token = apikey.Token
	}
	fmt.Println(apikey.ID)
	b.AppKeyId = int(apikey.ID)
	return nil
}

func (b *BaiduWX) getToken(key string, sec string) (token *BDTokenResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://wenxin.baidu.com/moduleApi/portal/api/oauth/token?grant_type=client_credentials&client_id=%s&client_secret=%s", key, sec), nil)
	if err != nil {
		return
	}
	req = req.WithContext(b.Ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = b.sendRequest(req, &token)
	return
}

func (b *BaiduWX) sendRequest(req *http.Request, v interface{}) error {
	res, err := b.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil || errRes.Error == nil {
			return fmt.Errorf("status code: %d", res.StatusCode)
		}
		return fmt.Errorf("status code: %d, message: %s", res.StatusCode, errRes.Error.Message)
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
	}

	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (b *BaiduWX) Pic2Pic(request *BDImageTaskRequest) (task BDImageTask, err error) {

	e := b.getClient(0)
	if e != nil {
		err = e
		return
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = writer.WriteField("text", request.Prompt)
	if err != nil {
		return
	}
	err = writer.WriteField("style", request.Prompt)
	if err != nil {
		return
	}
	err = writer.WriteField("num", "1")
	if err != nil {
		return
	}
	err = writer.WriteField("resolution", "1024*1024")
	if err != nil {
		return
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", request.Image.Filename))
	h.Set("Content-Type", "image/jpeg")
	img, err := writer.CreatePart(h)
	if err != nil {
		return
	}

	f, err := request.Image.Open()
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(img, f)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://wenxin.baidu.com/moduleApi/portal/api/rest/1.0/ernievilg/v1/txt2img?access_token=%s", b.Token), body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	taskResponse := &BDImageTaskResponse{}
	err = b.sendRequest(req, &taskResponse)
	if taskResponse.Code != 0 {
		if taskResponse.Msg == "暂不支持创作该内容，请修改后再试" {
			taskResponse.Msg = "内容不合规，请重新选择图片再试"
		}
		err = errors.New(taskResponse.Msg)
		return
	}
	task = BDImageTask{
		TaskId:   taskResponse.Data.TaskId,
		AppKeyId: b.AppKeyId,
	}
	return

}

func (b *BaiduWX) Pic2PicTask(taskId string, appKeyId int) (res BDImageTaskResult, err error) {

	e := b.getClient(appKeyId)
	if e != nil {
		err = e
		return
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://wenxin.baidu.com/moduleApi/portal/api/rest/1.0/ernievilg/v1/getImg?access_token=%s&taskId=%s", b.Token, taskId), nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	taskResponse := &BDImageTaskResultResponse{}
	err = b.sendRequest(req, &taskResponse)
	logx.Info(taskResponse)
	if taskResponse.Code != 0 {
		err = errors.New(taskResponse.Msg)
		return
	}
	res = taskResponse.Data
	return

}
