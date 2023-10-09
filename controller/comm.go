package controller

// Response 基础响应结构体
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
