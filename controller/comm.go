package controller

// Response 基础响应结构体
type Response struct {
	Code int32  `json:"status_code,omitempty"`
	Msg  string `json:"status_msg,omitempty"`
}
