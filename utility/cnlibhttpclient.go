package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	Message string `json:"Message"`
}

type CNLibHTTPClient struct {
	token string
}

func (client *CNLibHTTPClient) Get(url string) (*Response, bool) {
	Logger.Debug("GET " + url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Error(err.Error())
	}
	return client.request(request)
}

func (client *CNLibHTTPClient) GetWithToken(url string) (*Response, bool) {
	Logger.Debug("GETWT " + url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Error(err.Error())
	}
	request.Header.Add("Token", client.token)
	return client.request(request)
}

func (client *CNLibHTTPClient) PostForm(url string, data url.Values) (*Response, bool) {
	Logger.Debug("POSTFORM " + url)
	request, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		Logger.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.request(request)
}

func (client *CNLibHTTPClient) request(request *http.Request) (*Response, bool) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		Logger.Error(err.Error())
	}
	if response.Header.Get("Token") != "" {
		client.token = response.Header.Get("Token")
	}
	return client.handleResponse(response)
}

func (client *CNLibHTTPClient) handleResponse(response *http.Response) (*Response, bool) {
	message := &Response{}
	ResponseToJSON(response, message)
	Logger.Debug(fmt.Sprintf("Get response with status code: %v with message: %s", response.StatusCode, message.Message))
	return message, response.StatusCode == http.StatusOK
}

func NewCNLibHTTPClient() *CNLibHTTPClient {
	return &CNLibHTTPClient{}
}

func ResponseToString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Logger.Error(err.Error())
	}
	defer response.Body.Close()
	return string(body)
}

func ResponseToJSON(response *http.Response, jsonStruct interface{}) {
	err := json.NewDecoder(response.Body).Decode(jsonStruct)
	if err != nil {
		Logger.Error(err.Error())
	}
	defer response.Body.Close()
}
