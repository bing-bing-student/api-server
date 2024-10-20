package utils

import (
	"blog/global"
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
	"time"
)

type shortMessage struct {
	Code string `json:"code"`
}

func createClient() (result *dysmsapi20170525.Client, err error) {
	configInfo := &openapi.Config{
		AccessKeyId:     tea.String(global.CONFIG.ALYConfig.AccessKeyId),
		AccessKeySecret: tea.String(global.CONFIG.ALYConfig.AccessKeySecret),
	}
	configInfo.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(configInfo)
	return result, err
}

func SendMessage(phone string) (err error) {
	client, err := createClient()
	if err != nil {
		return err
	}

	code, err := GenerateRandomDigits(6)
	if err != nil {
		return err
	}
	// 将code存入redis当中
	if _, err := global.RedisSentinel.Set(global.CONTEXT, "code", code, time.Minute).Result(); err != nil {
		return err
	}

	identifyingCode := shortMessage{Code: code}
	identifyingCodeString, err := json.Marshal(identifyingCode)
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("BingBingBlog"),
		TemplateCode:  tea.String("SMS_296325152"),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(string(identifyingCodeString)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (e error) {
		defer func() {
			if err = tea.Recover(recover()); err != nil {
				e = err
			}
		}()
		if _, err := client.SendSmsWithOptions(sendSmsRequest, runtime); err != nil {
			return err
		}

		return err
	}()

	if tryErr != nil {
		var e = &tea.SDKError{}
		var t *tea.SDKError
		if errors.As(tryErr, &t) {
			e = t
		}
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(e.Data)))
		if err = d.Decode(&data); err != nil {
			return err
		}
		if _, err = util.AssertAsString(e.Message); err != nil {
			return err
		}
	}
	return err
}
