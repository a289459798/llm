package utils

import (
	"chatgpt-tools/model"
	"context"
	"errors"
	"fmt"
	"github.com/importcjj/sensitive"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"gorm.io/gorm"
)

func Filter(str string, db *gorm.DB) string {
	filter := sensitive.New()
	filter.LoadWordDict("data/sensitive_words_lines.txt")
	valid, s := filter.Validate(str)
	if valid {
		return ""
	}
	db.Create(&model.Contraband{
		Content: str,
		Error:   s,
	})
	return fmt.Sprintf("违禁词：%s，请修改后提交", s)
}

func ImageCheck(ctx context.Context, ak string, sk string, bucket string, key string) (bool, error) {
	a := auth.New(ak, sk)
	c := client.DefaultClient
	body := map[string]interface{}{
		"data": map[string]string{
			"uri": fmt.Sprintf("qiniu:///%s/%s", bucket, key),
		},
		"params": map[string][]string{
			"scenes": {
				"pulp",
				"terror",
				"politician",
			},
		},
	}

	var ret map[string]interface{}
	err := c.CredentialedCallWithJson(ctx, a, auth.TokenQiniu, &ret, "POST", "http://ai.qiniuapi.com/v3/image/censor", nil, body)
	if err != nil {
		return false, err
	}

	if ret["code"].(float64) != 200 {
		return false, errors.New(ret["message"].(string))
	}
	if ret["result"].(map[string]interface{})["suggestion"] != "pass" {
		return false, errors.New("图片涉嫌违规")
	}
	return true, nil
}
