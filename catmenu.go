package yapi

import (
	"encoding/json"
	"github.com/jinzhu/copier"
)

type CatMenuData []CatData

type CatData struct {
	ID   int    `json:"_id" structs:"_id"`
	UID  int    `json:"uid" structs:"uid"`
	Name string `json:"name" structs:"name"`
	Desc string `json:"desc" structs:"desc"`
}

type CatMenu struct {
	ErrCode int         `json:"errcode" structs:"errcode"`
	ErrMsg  string      `json:"errmsg" structs:"errmsg"`
	Data    CatMenuData `json:"data" structs:"data"`
	string  string      `json:"string"`
}

func (p *CatMenu) ToString() string {
	return p.string
}

// CatMenuService .
type CatMenuService struct {
	client *Client
}

type CatMenuParam struct {
	Token     string `url:"token"`
	ProjectID int    `url:"project_id"`
}

type ModifyMenuParam struct {
	ProjectID int `json:"project_id" url:"project_id"`
	CatData
}

type ModifyMenumReq struct {
	Token string `json:"token"`
	ModifyMenuParam
}

type ModifyMenuResp struct {
	CommonResp
	Data   interface{} `json:"data" structs:"data"`
	string string      `json:"string"`
}

func (p *ModifyMenuResp) ToString() string {
	return p.string
}

func (s *CatMenuService) Get(projectId int) (*CatMenu, error) {
	apiEndpoint := "api/interface/getCatMenu"
	catMenuParam := CatMenuParam{}
	catMenuParam.ProjectID = projectId
	catMenuParam.Token = s.client.Authentication.token
	url, err := addOptions(apiEndpoint, &catMenuParam)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	result := CatMenu{}
	result.string = resp
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}

func (s *CatMenuService) AddOrUpdate(param *ModifyMenuParam) (*ModifyMenuResp, error) {
	apiEndpoint := "api/interface/add_cat"
	modifyMenuReq := ModifyMenumReq{}
	modifyMenuReq.Token = s.client.Authentication.token
	copier.Copy(&modifyMenuReq, param)

	resp, err := s.client.Post(apiEndpoint, modifyMenuReq)
	if err != nil {
		return nil, err
	}
	result := ModifyMenuResp{}
	result.string = resp
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}
