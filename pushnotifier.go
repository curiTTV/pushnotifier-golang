package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// New creates a new instance of Pushnotifier
func New(packageName, apiToken string) (*Pushnotifier, error) {
	packageName = strings.TrimSpace(packageName)
	apiToken = strings.TrimSpace(apiToken)

	if packageName == "" {
		return nil, ErrMissingPackageName
	} else if apiToken == "" {
		return nil, ErrMissingAPIToken
	}

	pn := &Pushnotifier{
		packageName: packageName,
		apiToken:    apiToken,
		endpoints:   make(map[string]*endpoint),
	}

	pn.endpoints["login"] = &endpoint{
		method: "POST",
		uri:    "user/login",
	}

	pn.endpoints["refreshToken"] = &endpoint{
		method: "GET",
		uri:    "user/refresh",
	}

	pn.endpoints["devices"] = &endpoint{
		method: "GET",
		uri:    "devices",
	}

	pn.endpoints["sendText"] = &endpoint{
		method: "PUT",
		uri:    "notifications/text",
	}

	pn.endpoints["sendURL"] = &endpoint{
		method: "PUT",
		uri:    "notifications/url",
	}

	pn.endpoints["sendNotification"] = &endpoint{
		method: "PUT",
		uri:    "notifications/notification",
	}

	pn.endpoints["sendImage"] = &endpoint{
		method: "PUT",
		uri:    "notifications/image",
	}

	return pn, nil
}

// Login retrieves an app_token from pushnotifier.de
func (pn *Pushnotifier) Login(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return ErrMissingCredentials
	}

	pn.credentials = &credentials{
		Username: username,
		Password: password,
	}

	statusCode, resp, err := pn.apiRequest(pn.endpoints["login"], pn.credentials)
	if err != nil {
		return err
	}

	if statusCode >= 500 {
		return ErrPushnotifierServerError
	} else if statusCode != 200 {
		return ErrInvalidCredentials
	}

	pn.appToken = &appToken{}

	err = json.Unmarshal(resp, &pn.appToken)
	if err != nil {
		return err
	}

	return nil
}

func (pn *Pushnotifier) Devices() ([]*Device, error) {
	_, resp, err := pn.apiRequest(pn.endpoints["devices"], nil)
	if err != nil {
		return nil, err
	}

	var devices []*Device

	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (pn *Pushnotifier) apiRequest(endpoint *endpoint, body interface{}) (int, []byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	url := "https://api.pushnotifier.de/v2/" + endpoint.uri

	var request *http.Request

	switch endpoint.method {
	case "GET":
		request, err = http.NewRequest("GET", url, nil)

	case "POST", "PUT":
		request, err = http.NewRequest(endpoint.method, url, bytes.NewBuffer(data))

	default:
		return 0, nil, ErrInvalidHTTPMethod
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(pn.packageName+":"+pn.apiToken)))
	if pn.appToken != nil && pn.appToken.Token != "" {
		request.Header.Set("X-AppToken", pn.appToken.Token)
	}

	if err != nil {
		return 0, nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, respBody, nil
}
