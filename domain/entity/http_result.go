package entity

type HttpResult struct {
	ID            int64  `json:"id"`             // id
	Time          int64  `json:"time"`           // 时间
	TaskName      string `json:"task_name"`      // 任务的名称
	Target        string `json:"target"`         // http地址
	TargetTo      string `json:"target_to"`      // http地址 实际地址
	StatusCode    int    `json:"status_code"`    // HTTP状态码
	ContentLength int64  `json:"content_length"` // HTTP响应内容大小
	ContentType   string `json:"content_type"`   // HTTP响应内容类型
	Raw           string `json:"raw"`            // HTTP响应内容
	Title         string `json:"title"`          // 网页标题
	Description   string `json:"description"`    // 网页摘要
}
