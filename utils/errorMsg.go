package utils

//1000-1999 用户模块错误
//2000-2999 文章模块错误
//3000-3999 分类模块错误

const (
	SUCCESS             = 200
	ERROR               = 500
	ErrorUsernameWrong  = 1001
	ErrorPasswordWrong  = 1002
	ErrorUserNotExist   = 1003
	ErrorTokenExist     = 1004
	ErrorTokenRuntime   = 1005
	ErrorTokenWrong     = 1006
	ErrorTokenTypeWrong = 1007
	ErrorUserNoRight    = 1008
	ErrorArtNotExist    = 2000
	ErrorCateUsed       = 3000
	ErrorCateNotExist   = 3002
)

var codeMsg = map[int]string{
	SUCCESS:             "OK",
	ERROR:               "FAIL",
	ErrorUsernameWrong:  "用户名错误",
	ErrorPasswordWrong:  "密码错误",
	ErrorUserNotExist:   "用户不存在",
	ErrorTokenExist:     "TOKEN不存在,请重新登陆",
	ErrorTokenRuntime:   "TOKEN已过期,请重新登陆",
	ErrorTokenWrong:     "TOKEN不正确,请重新登陆",
	ErrorTokenTypeWrong: "TOKEN格式错误,请重新登陆",
	ErrorUserNoRight:    "该用户无权限",
	ErrorArtNotExist:    "文章不存在",
	ErrorCateUsed:       "该分类已存在",
	ErrorCateNotExist:   "该分类不存在",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
