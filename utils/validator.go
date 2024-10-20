package utils

import (
	"blog/global"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

func Translate(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		value.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是支持的语言环境（支持多个）
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		global.Translate, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(value, global.Translate)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(value, global.Translate)
		default:
			err = enTranslations.RegisterDefaultTranslations(value, global.Translate)
		}
		return
	}
	return
}

// CustomRules 定义一个自定义的校验规则
func CustomRules(tag string) (err error) {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 在校验器注册自定义的校验方法
		if err := value.RegisterValidation(tag, customFunc); err != nil {
			return err
		}
	}
	return err
}

func customFunc(fl validator.FieldLevel) bool {
	// 正则表达式来匹配字母、数字和下划线，且第一个字符不能是下划线
	re := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*$`)
	return re.MatchString(fl.Field().String())
}
