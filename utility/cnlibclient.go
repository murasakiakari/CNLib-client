package utility

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"
)

type CNLibConfiguration struct {
	BackendURL string `yaml:"backendURL"`
	DebugMode  bool   `yaml:"debugMode"`
}

func getConfiguration() CNLibConfiguration {
	configurationFile := CurrentWorkingDirectory.Join("configuration.yml")
	configuration := &CNLibConfiguration{}
	if configurationFile.IsExist() {
		fileData, err := configurationFile.ReadFile()
		if err != nil {
			Logger.Error("Cannot read configuration file at " + configurationFile)
		}
		err = yaml.Unmarshal(fileData, configuration)
		if err != nil {
			Logger.Error("Cannot unmarshal configuration file at " + configurationFile)
		}
	} else {
		Logger.Error("Configuration file " + configurationFile + " does not exist")
	}
	return *configuration
}

type CNLibClient struct {
	configuration CNLibConfiguration
	httpClient    *CNLibHTTPClient
}

func (client *CNLibClient) Run() {
	fmt.Println("Welcome to CNLib-client")
	client.printHelp()
	Loop: for {
		input := Input("Please enter the function you want: ")
		switch input {
		case "1":
			response, _ := client.login()
			fmt.Println(response.Message)
		case "2":
			response, _ := client.logout()
			fmt.Println(response.Message)
		case "3":
			response, _ := client.test()
			fmt.Println(response.Message)
		case "5":
			break Loop
		default:
			client.printHelp()
		}
	}
}

func (client *CNLibClient) printHelp() {
	fmt.Println("Function: ")
	fmt.Println("1. login")
	fmt.Println("2. logout")
	fmt.Println("3. test")
	fmt.Println("4. print help message")
	fmt.Println("5. quit")
}

func (client *CNLibClient) login() (*Response, bool) {
	account := Input("Please input your account: ")
	password := Input("Please input your password: ")
	return client.httpClient.PostForm(client.configuration.BackendURL + "/login", url.Values{"account": {account}, "password": {password}})
}

func (client *CNLibClient) logout() (*Response, bool) {
	return client.httpClient.GetWithToken(client.configuration.BackendURL + "/logout")
}

func (client *CNLibClient) test() (*Response, bool) {
	return client.httpClient.GetWithToken(client.configuration.BackendURL + "/test")
}

func NewCNLibClient() *CNLibClient {
	client := &CNLibClient{configuration: getConfiguration(), httpClient: NewCNLibHTTPClient()}
	Logger.SetDebugMode(client.configuration.DebugMode)
	return client
}
