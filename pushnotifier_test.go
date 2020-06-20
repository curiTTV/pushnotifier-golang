package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type (
	mockClient struct {
		mockDo func(req *http.Request) (*http.Response, error)
	}

	testData struct {
		username    string
		password    string
		packageName string
		apiToken    string
		appToken    string
	}
)

var (
	pn   *Pushnotifier
	data *testData
)

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	return c.mockDo(req)
}

func init() {
	data = &testData{
		username:    "testuser",
		password:    "testpass",
		packageName: "testpackage",
		apiToken:    "apitoken123",
		appToken:    "apptoken123",
	}
	pn = New(data.packageName, data.apiToken)
}

func TestLogin(t *testing.T) {
	wantBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(data.packageName+":"+data.apiToken))
	responseData := "{\"username\": \"" + data.username + "\",\"avatar\": \"\",\"app_token\": \"" + data.appToken + "\",\"expires_at\": 1513637432}"

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Error("Invalid method for devices provided")
			}

			if req.Header.Get("Authorization") != wantBasicAuth {
				t.Error("Authorization header invalid")
			}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Error(err)
			}

			cred := &credentials{}
			err = json.Unmarshal(body, cred)
			if err != nil {
				t.Error(err)
			}

			if cred.Username != data.username || cred.Password != data.password {
				t.Error("username or password not correctly provided")
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseData))),
			}, nil
		},
	}

	err := pn.Login(data.username, data.password)
	if err != nil {
		t.Error(err)
	} else if pn.appToken.Token != data.appToken {
		t.Error("Invalid apptoken in Pushnotifier struct")
	}
}

func TestRefresh(t *testing.T) {
	wantAppToken := &appToken{
		Token:     "abc123",
		ExpiresAt: 123,
	}

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Error("Invalid method for refresh provided")
			}

			appToken, err := json.Marshal(wantAppToken)
			if err != nil {
				t.Error(err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(appToken)),
			}, nil
		},
	}

	if err := pn.Refresh(); err != nil {
		t.Error(err)
	}
}

func TestDevices(t *testing.T) {
	wantDevices := []*Device{
		{
			ID:    "1",
			Title: "t1",
			Model: "m1",
			Image: "i1",
		},
		{
			ID:    "2",
			Title: "t2",
			Model: "m2",
			Image: "i2",
		},
	}

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Error("Invalid method for devices provided")
			}

			devices, err := json.Marshal(wantDevices)
			if err != nil {
				t.Error(err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(devices)),
			}, nil
		},
	}

	gotDevices, err := pn.Devices()
	if err != nil {
		t.Error(err)
	}

	for i, device := range gotDevices {
		if device.ID != wantDevices[i].ID {
			t.Error("Invalid id returned")
		}
		if device.Title != wantDevices[i].Title {
			t.Error("Invalid title returned")
		}
		if device.Model != wantDevices[i].Model {
			t.Error("Invalid model returned")
		}
		if device.Image != wantDevices[i].Image {
			t.Error("Invalid image returned")
		}
	}
}

func TestText(t *testing.T) {
	wantDataAll := []*Device{
		{
			ID: "d1",
		},
		{
			ID: "d2",
		},
	}
	wantDataSpecific := []*Device{
		{
			ID: "d3",
		},
	}

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Error("Invalid method for test notifications provided")
			}

			response := &pnNotificationResponse{}

			for _, deviceID := range wantDataAll {
				response.Success = append(response.Success, deviceID.ID)
			}

			data, err := json.Marshal(response)
			if err != nil {
				t.Error(err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(data)),
			}, nil
		},
	}

	if err := pn.Text(nil, "hello"); err != nil {
		t.Error(err)
	}

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Error("Invalid method for test notifications provided")
			}

			response := &pnNotificationResponse{}

			for _, deviceID := range wantDataSpecific {
				response.Success = append(response.Success, deviceID.ID)
			}

			data, err := json.Marshal(response)
			if err != nil {
				t.Error(err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(data)),
			}, nil
		},
	}

	if err := pn.Text(wantDataSpecific, "hello"); err != nil {
		t.Error(err)
	}
}
