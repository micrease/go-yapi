package yapi

/**
接口通用返回
*/
type CommonResp struct {
	ErrCode int    `json:"errcode" structs:"errcode"`
	ErrMsg  string `json:"errmsg" structs:"errmsg"`
}

type ModifyResult struct {
	Ok        int `json:"ok" structs:"ok"`
	NModified int `json:"nModified" structs:"nModified"`
	N         int `json:"n" structs:"n"`
}

type ModifyResp struct {
	CommonResp
	Data   interface{} `json:"data" structs:"data"`
	string string      `json:"string"`
}

func (m *ModifyResp) ToString() string {
	return m.string
}
