package yapi

import (
	"encoding/json"
	"html/template"
)

// InterfaceService .
type InterfaceService struct {
	client *Client
}

type ReqKVItemSimple struct {
	Name    string `json:"name" structs:"name"`
	Value   string `json:"value" structs:"value"`
	Example string `json:"example" structs:"example"`
	Desc    string `json:"desc" structs:"desc"`
}

type ReqKVItemDetail struct {
	ReqKVItemSimple
	Type     string `json:"type" structs:"type"`
	Required string `json:"required" structs:"required"`
}

type interfaceBase struct {
	ID        int      `json:"_id" structs:"_id"`
	UID       int      `json:"uid" structs:"uid"`
	CatID     int      `json:"catid" structs:"catid"`
	ProjectID int      `json:"project_id" structs:"project_id"`
	EditUID   int      `json:"edit_uid" structs:"edit_uid"`
	AddTime   int      `json:"add_time" structs:"add_time"`
	UpTime    int      `json:"up_time" structs:"up_time"`
	Status    string   `json:"status" structs:"status"`
	Title     string   `json:"title" structs:"title"`
	Path      string   `json:"path" structs:"path"`
	Method    string   `json:"method" structs:"method"`
	Tag       []string `json:"tag" structs:"tag"`
}

type interfaceReq struct {
	ReqParams           []ReqKVItemSimple `json:"req_params" structs:"req_params"`
	ReqHeaders          []ReqKVItemDetail `json:"req_headers" structs:"req_headers"`
	ReqQuery            []ReqKVItemDetail `json:"req_query" structs:"req_query"`
	ReqBodyForm         []ReqKVItemDetail `json:"req_body_form" structs:"req_body_form"`
	ReqBodyIsJsonSchema bool              `json:"req_body_is_json_schema" structs:"req_body_is_json_schema"`
	ReqBodyType         string            `json:"req_body_type" structs:"req_body_type"`
	ReqBodyOther        template.HTML     `json:"req_body_other" structs:"req_body_other"`
}

type interfaceRes struct {
	ResBodyIsJsonSchema bool          `json:"res_body_is_json_schema" structs:"res_body_is_json_schema"`
	ResBodyType         string        `json:"res_body_type" structs:"res_body_type"`
	ResBody             template.HTML `json:"res_body" structs:"res_body"`
}

type InterfaceData struct {
	interfaceBase
	interfaceReq
	interfaceRes
}

type AddOrUpdateInterfaceData struct {
	Token string `json:"token" structs:"token"`
	InterfaceData
}

type Interface struct {
	CommonResp
	Data   InterfaceData `json:"data" structs:"data"`
	string string        `json:"string"`
}

func (i *Interface) ToString() string {
	return i.string
}

type InterfaceParam struct {
	Token string `url:"token"`
	ID    int    `url:"id"`
}

type InterfaceListData struct {
	Count int             `json:"count" structs:"count"`
	Total int             `json:"total" structs:"total"`
	List  []InterfaceData `json:"list" structs:"list"`
}

type InterfaceList struct {
	ErrCode int               `json:"errcode" structs:"errcode"`
	ErrMsg  string            `json:"errmsg" structs:"errmsg"`
	Data    InterfaceListData `json:"data" structs:"data"`
	string  string            `json:"string"`
}

func (i *InterfaceList) ToString() string {
	return i.string
}

type InterfaceListParam struct {
	Token string `url:"token,omitempty"`
	CatID int    `url:"catid,omitempty"`
	Page  int    `url:"Page"`
	Limit int    `url:"limit"`
}

type UploadSwaggerReq struct {
	Type  string `json:"type" structs:"types"`
	Json  string `json:"json" structs:"jsons"`
	Merge string `json:"merge" structs:"string"`
	Token string `json:"token" structs:"token"`
	url   string `json:"url" structs:"url"`
}

func (s *InterfaceService) GetList(opt *InterfaceListParam) (*InterfaceList, error) {
	apiEndpoint := "api/interface/list_cat"
	opt.Token = s.client.Authentication.token
	url, err := addOptions(apiEndpoint, opt)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	result := InterfaceList{}
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}

func (s *InterfaceService) Get(id int) (*Interface, error) {
	apiEndpoint := "api/interface/get"
	interfaceParam := InterfaceParam{}
	interfaceParam.ID = id
	interfaceParam.Token = s.client.Authentication.token
	url, err := addOptions(apiEndpoint, &interfaceParam)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	result := Interface{}
	result.string = resp
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}

func (s *InterfaceService) AddOrUpdate(data *InterfaceData) (*ModifyResp, error) {
	apiEndpoint := "api/interface/save"
	addOrUpdateInterfaceData := AddOrUpdateInterfaceData{}
	addOrUpdateInterfaceData.Token = s.client.Authentication.token
	resp, err := s.client.Post(apiEndpoint, addOrUpdateInterfaceData)
	if err != nil {
		return nil, err
	}
	result := ModifyResp{}
	result.string = resp
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}

func (s *InterfaceService) UploadSwagger(data *string) (*ModifyResp, error) {
	apiEndpoint := "api/open/import_data"
	uploadSwaggerReq := new(UploadSwaggerReq)
	uploadSwaggerReq.Token = s.client.Authentication.token
	uploadSwaggerReq.Type = "swagger"
	uploadSwaggerReq.Merge = "merge"
	uploadSwaggerReq.Json = *data

	resp, err := s.client.Post(apiEndpoint, uploadSwaggerReq)
	if err != nil {
		return nil, err
	}
	result := ModifyResp{}
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}
