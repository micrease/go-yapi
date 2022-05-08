package yapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	testInstanceURL = "http://yapi.iacorn.cn"
	testToken = "8cde6e3bcbdb1c7bc0d8a2992fe0d79b0f07d186bfe7ebc5c6723b403e6830d8"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the client being tested.
	testClient *Client

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

type testValues map[string]string

type BodyTest struct {
	ID string `json:"id,omitempty" structs:"id,omitempty"`
}

// setup sets up a test HTTP server along with a Client that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup() {
	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// test client configured to use test server
	testClient, _ = NewClient(nil, testServer.URL, "")
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL.String(); !strings.HasPrefix(got, want) {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}

func testRequestParams(t *testing.T, r *http.Request, want map[string]string) {
	params := r.URL.Query()

	if len(params) != len(want) {
		t.Errorf("Request params: %d, want %d", len(params), len(want))
	}

	for key, val := range want {
		if got := params.Get(key); val != got {
			t.Errorf("Request params: %s, want %s", got, val)
		}
	}
}

func TestNewClient_WrongUrl(t *testing.T) {
	c, err := NewClient(nil, ":/3000/", "")

	if err == nil {
		t.Error("Expected an error. Got none")
	}
	if c != nil {
		t.Errorf("Expected no client. Got %+v", c)
	}
}

func TestNewClient_WithHttpClient(t *testing.T) {
	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Minute
	c, err := NewClient(httpClient, testInstanceURL, "")

	if err != nil {
		t.Errorf("Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
	}
	if !reflect.DeepEqual(c.client, httpClient) {
		t.Errorf("HTTP clients are not equal. Injected %+v, got %+v", httpClient, c.client)
	}
}

func TestNewClient_WithServices(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, "")

	if err != nil {
		t.Errorf("Got an error: %s", err)
	}
	if c.Interface == nil {
		t.Error("No InterfaceService provided")
	}
	if c.Project == nil {
		t.Error("No ProjectService provided")
	}

	if c.Authentication == nil {
		t.Error("No AuthenticationService provided")
	}
}

func TestCheckResponse(t *testing.T) {
	codes := []int{
		http.StatusOK, http.StatusPartialContent, 299,
	}

	for _, c := range codes {
		r := &http.Response{
			StatusCode: c,
		}
		if err := CheckResponse(r); err != nil {
			t.Errorf("CheckResponse throws an error: %s", err)
		}
	}
}

func TestClient_NewRequest(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, "")
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "api/", testInstanceURL+"api/"
	inBody, outBody := &BodyTest{ID: "1"}, `{"id":"1"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}
}

func TestClient_NewRawRequest(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, "")
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "api/", testInstanceURL+"api/"

	outBody := `{"id":"1"}` + "\n"
	inBody := outBody
	req, _ := c.NewRawRequest("GET", inURL, strings.NewReader(outBody))

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRawRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRawRequest(%v) Body is %v, want %v", inBody, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestClient_NewRequest_BadURL(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, "")
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}
	_, err = c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

// If a nil body is passed to client.NewRequest, make sure that nil is also passed to http.NewRequest.
// In most cases, passing an io.Reader that returns no content is fine,
// since there is no difference between an HTTP request body that is an empty string versus one that is not set at all.
// However in certain cases, intermediate systems may treat these differently resulting in subtle errors.
func TestClient_NewRequest_EmptyBody(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, "")
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	body := new(foo)
	_, err := testClient.Do(req, body)
	if err != nil {
		t.Error("Expected error to be returned.")
	}

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestClient_Do_HTTPResponse(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	res, _ := testClient.Do(req, nil)
	_, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error on parsing HTTP Response = %v", err.Error())
	} else if res.StatusCode != 200 {
		t.Errorf("Response code = %v, want %v", res.StatusCode, 200)
	}
}

func TestClient_Do_HTTPError(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do() function.
// A redirect loop is pretty unlikely to occur within the API, but does allow us to exercise the right code path.
func TestClient_Do_RedirectLoop(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %+v.", err)
	}
}

func TestClient_GetBaseURL_WithURL(t *testing.T) {
	u, err := url.Parse(testInstanceURL)
	if err != nil {
		t.Errorf("URL parsing -> Got an error: %s", err)
	}

	c, err := NewClient(nil, testInstanceURL, "")
	if err != nil {
		t.Errorf("Client creation -> Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
	}

	if b := c.GetBaseURL(); !reflect.DeepEqual(b, *u) {
		t.Errorf("Base URLs are not equal. Expected %+v, got %+v", *u, b)
	}
}

func TestClient_AddOrUpdateInterfaceData(t *testing.T) {
	c, err := NewClient(nil, testInstanceURL, testToken)
	if err != nil {
		t.Errorf("Client creation -> Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
	}

	markMenu := "test"

	// 获取项目id
	project, _, _ := c.Project.Get()
	projectId := project.Data.ID
	fmt.Printf("项目id:%d\n", projectId)

	// 获取项目下的分类
	catMenu,_,_ := c.CatMenu.Get(projectId)
	fmt.Println("项目目录分类：")
	printResult(catMenu)

	var menuExsit *CatData
	for _, menu := range catMenu.Data {
		if menu.Name == markMenu {
			menuExsit = &menu
		}
	}

	if menuExsit == nil {
		// 不存在，创建新目录
		modifyMenuParam := new(ModifyMenuParam)
		modifyMenuParam.ProjectID = projectId
		modifyMenuParam.Name = markMenu
		modifyMenuResp,_,_ := c.CatMenu.AddOrUpdate(modifyMenuParam)
		fmt.Println(modifyMenuResp)
	}

	interfaceDataResp, _, _  := c.Interface.Get(415)
	printResult(interfaceDataResp)

	// 新增或者修改
	interfaceData := interfaceDataResp.Data
	interfaceData.Title = "测试修idea1111111改1111"
	interfaceData.ReqBodyType = "json"

	// 添加body
	reqKVItemSimple := new(ReqKVItemSimple)
	reqKVItemSimple.Name = "test"
	reqKVItemSimple.Desc = "测试字段"
	reqKVItemSimple.Example = "test"
	reqKVItemSimple.Value = ""

	reqKVItemDetail := new(ReqKVItemDetail)
	reqKVItemDetail.Type = "text"
	reqKVItemDetail.Value = ""
	reqKVItemDetail.Example = "111"
	reqKVItemDetail.Desc = "111"
	reqKVItemDetail.Required = "1"
	reqKVItemDetail.Name = "11111"

	details := append([]ReqKVItemDetail{}, *reqKVItemDetail)

	//marshal, err := json.Marshal(reqKVItemSimple)
	interfaceData.ReqBodyOther = "{\n    \"$schema\": \"http://json-schema.org/schema#\",\n    \"type\": \"object\",\n    \"properties\": {\n        \"Bar\": {\n            \"type\": \"string\"\n        },\n        \"Baz\": {\n            \"type\": \"array\",\n            \"items\": {\n                \"type\": \"string\"\n            }\n        },\n        \"List\": {\n            \"type\": \"array\",\n            \"items\": {\n                \"type\": \"object\",\n                \"properties\": {\n                    \"Value\": {\n                        \"type\": \"string\"\n                    }\n                },\n                \"required\": [\n                    \"Value\"\n                ]\n            }\n        },\n        \"Qux\": {\n            \"type\": \"integer\"\n        },\n        \"Zoo\": {\n            \"type\": \"string\"\n        },\n        \"foo\": {\n            \"type\": \"boolean\"\n        }\n    },\n    \"required\": [\n        \"foo\",\n        \"Qux\",\n        \"Baz\",\n        \"Zoo\",\n        \"List\"\n    ]\n}"
	interfaceData.ID = 0
	interfaceData.ProjectID = projectId
	interfaceData.CatID = menuExsit.ID

	//interfaceData.ResBodyIsJsonSchema = false
	interfaceData.ReqBodyForm = details
	interfaceData.ReqParams = append([]ReqKVItemSimple{}, *reqKVItemSimple)
	interfaceData.ReqQuery= details
	interfaceData.ReqHeaders = details

	addOrUpdateResp, _, _:= c.Interface.AddOrUpdate(&interfaceData)
	printResult(addOrUpdateResp)
}

func printResult(result interface{}) {
	marshal,_ := json.Marshal(result)
	api := string(marshal)
	fmt.Println(api)
}

func TestClient_UploadSwagger(t *testing.T) {
	// 使用ioutil一次性读取文件
	data, err := ioutil.ReadFile("E:\\golang\\global_gopath\\src\\github.com\\swaggo\\swag\\testdata\\composition\\docs\\swagger.json")
	if err != nil {
		fmt.Println("read file err:", err.Error())
		return
	}

	c, err := NewClient(nil, testInstanceURL, testToken)
	if err != nil {
		t.Errorf("Client creation -> Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
	}

	swagger := string(data)
	result,_,_ := c.Interface.UploadSwagger(&swagger)
	printResult(result)
}