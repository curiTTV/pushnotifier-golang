package pushnotifier

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
	wantDevices := []*Device{
		{
			ID: "d1",
		},
		{
			ID: "d2",
		},
	}
	wantContent := "hello"

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Error("Invalid method for test notifications provided")
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			notification := &pnNotification{}
			err := json.Unmarshal(buf.Bytes(), notification)
			if err != nil {
				t.Error(err)
			}

			for i, deviceID := range notification.Devices {
				if deviceID != wantDevices[i].ID {
					t.Error("mismatch of requested device id")
				}
			}

			if notification.Content != wantContent {
				t.Error("content mismatch")
			}

			response := &pnNotificationResponse{}

			for _, deviceID := range wantDevices {
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

	// all devices
	pn.devices = wantDevices

	if err := pn.Text(nil, wantContent); err != nil {
		t.Error(err)
	}

	// specific devices
	wantDevices = []*Device{
		{
			ID: "d3",
		},
	}

	if err := pn.Text(wantDevices, wantContent); err != nil {
		t.Error(err)
	}
}

func TestURL(t *testing.T) {
	wantDevices := []*Device{
		{
			ID: "d1",
		},
		{
			ID: "d2",
		},
	}
	wantURL := "https://github.com/curiTTV"

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Error("Invalid method for test notifications provided")
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			notification := &pnNotification{}
			err := json.Unmarshal(buf.Bytes(), notification)
			if err != nil {
				t.Error(err)
			}

			for i, deviceID := range notification.Devices {
				if deviceID != wantDevices[i].ID {
					t.Error("mismatch of requested device id")
				}
			}

			if notification.URL != wantURL {
				t.Error("url mismatch")
			}

			response := &pnNotificationResponse{}

			for _, deviceID := range wantDevices {
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

	wantDevices = []*Device{
		{
			ID: "d5",
		},
	}

	if err := pn.URL(wantDevices, wantURL); err != nil {
		t.Error(err)
	}
}

func TestNotification(t *testing.T) {
	wantDevices := []*Device{
		{
			ID: "d1",
		},
		{
			ID: "d2",
		},
	}
	wantURL := "https://github.com/curiTTV"
	wantContent := "hello"

	httpClient = &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Error("Invalid method for test notifications provided")
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			notification := &pnNotification{}
			err := json.Unmarshal(buf.Bytes(), notification)
			if err != nil {
				t.Error(err)
			}

			for i, deviceID := range notification.Devices {
				if deviceID != wantDevices[i].ID {
					t.Error("mismatch of requested device id")
				}
			}

			if notification.URL != wantURL {
				t.Error("url mismatch")
			}

			if notification.Content != wantContent {
				t.Error("url mismatch")
			}

			response := &pnNotificationResponse{}

			for _, deviceID := range wantDevices {
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

	wantDevices = []*Device{
		{
			ID: "d5",
		},
	}

	if err := pn.Notification(wantDevices, wantContent, wantURL); err != nil {
		t.Error(err)
	}
}
