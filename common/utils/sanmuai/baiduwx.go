package sanmuai

import (
	"bytes"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type BaiduWX struct {
	Ctx        context.Context
	SvcCtx     *svc.ServiceContext
	HTTPClient *http.Client
	Token      string
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
	TaskId int `json:"taskId"`
}

func NewBaiduWX(c context.Context, svcCtx *svc.ServiceContext) *BaiduWX {
	return &BaiduWX{
		Ctx:        c,
		SvcCtx:     svcCtx,
		Token:      "",
		HTTPClient: &http.Client{},
	}
}

func (b *BaiduWX) getClient() error {
	apikey := &model.Apikey{}
	b.SvcCtx.Db.Where("channel = ?", "baidu").Where("status = ?", 1).Order("rand()").Limit(1).Find(apikey)

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
		b.SvcCtx.Db.Save(&apikey)
	} else {
		b.Token = apikey.Token
	}
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
	fmt.Println(req.Header.Get("Content-Type"))
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

func (b *BaiduWX) Pic2Pic(request *BDImageTaskRequest) (task BDImageTask, err error) {

	e := b.getClient()
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

	img, err := writer.CreateFormFile("image", request.Image.Filename)
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
	fmt.Println(taskResponse)
	if taskResponse.Code != 0 {
		err = errors.New(taskResponse.Msg)
		return
	}
	task = BDImageTask{
		TaskId: taskResponse.Data.TaskId,
	}
	return

}
