package yapi

import (
	"encoding/json"
)

type EnvHeader struct {
	ID    string `json:"_id,omitempty" structs:"_id,omitempty"`
	Name  string `json:"name,omitempty" structs:"name,omitempty"`
	Value string `json:"value,omitempty" structs:"value,omitempty"`
}

type EnvGlobal struct {
	ID    string `json:"_id,omitempty" structs:"_id,omitempty"`
	Name  string `json:"name,omitempty" structs:"name,omitempty"`
	Value string `json:"value,omitempty" structs:"value,omitempty"`
}

type ProjectEnv struct {
	Header []EnvHeader `json:"header,omitempty" structs:"header,omitempty"`
	Global []EnvGlobal `json:"global,omitempty" structs:"global,omitempty"`
	ID     string      `json:"_id" structs:"_id"`
	Name   string      `json:"name" structs:"name"`
	Domain string      `json:"domain" structs:"domain"`
}

type ProjectData struct {
	ID      int          `json:"_id" structs:"_id"`
	UID     int          `json:"uid" structs:"uid"`
	GroupID int          `json:"group_id" structs:"group_id"`
	Name    string       `json:"name" structs:"name"`
	Role    bool         `json:"role" structs:"role"`
	Env     []ProjectEnv `json:"env" structs:"env"`
}

type Project struct {
	ErrCode int         `json:"errcode" structs:"errcode"`
	ErrMsg  string      `json:"errmsg" structs:"errmsg"`
	Data    ProjectData `json:"data" structs:"data"`
	string  string      `json:"string"`
}

func (p *Project) ToString() string {
	return p.string
}

// ProjectService .
type ProjectService struct {
	client *Client
}

type ProjectParam struct {
	Token string `url:"token"`
}

func (s *ProjectService) Get() (*Project, error) {
	apiEndpoint := "api/project/get"
	projectParam := ProjectParam{}
	projectParam.Token = s.client.Authentication.token
	url, err := addOptions(apiEndpoint, &projectParam)
	if err != nil {
		return nil, err
	}

	result := Project{}
	resp, err := s.client.Get(url, nil)
	if err != nil {
		return nil, err
	}

	result.string = resp
	err = json.Unmarshal([]byte(resp), &result)
	return &result, err
}
