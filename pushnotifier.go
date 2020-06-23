package pushnotifier

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func init() {
	httpClient = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
}

// New creates a new instance of Pushnotifier
func New(packageName, apiToken string) *Pushnotifier {
	packageName = strings.TrimSpace(packageName)
	apiToken = strings.TrimSpace(apiToken)

	pn := &Pushnotifier{
		packageName: packageName,
		apiToken:    apiToken,
	}

	pn.endpoints.login = &endpoint{
		method: "POST",
		uri:    "user/login",
	}

	pn.endpoints.refresh = &endpoint{
		method: "GET",
		uri:    "user/refresh",
	}

	pn.endpoints.devices = &endpoint{
		method: "GET",
		uri:    "devices",
	}

	pn.endpoints.sendText = &endpoint{
		method: "PUT",
		uri:    "notifications/text",
	}

	pn.endpoints.sendURL = &endpoint{
		method: "PUT",
		uri:    "notifications/url",
	}

	pn.endpoints.sendNotification = &endpoint{
		method: "PUT",
		uri:    "notifications/notification",
	}

	pn.endpoints.sendImage = &endpoint{
		method: "PUT",
		uri:    "notifications/image",
	}

	return pn
}

// Login retrieves an app_token from pushnotifier.de
func (pn *Pushnotifier) Login(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if pn.packageName == "" {
		return ErrMissingPackageName
	} else if pn.apiToken == "" {
		return ErrMissingAPIToken
	} else if username == "" || password == "" {
		return ErrMissingCredentials
	}

	pn.credentials = &credentials{
		Username: username,
		Password: password,
	}

	return pn.loginOrRefresh(pn.endpoints.login)
}

// Refresh refreshes the expired app token from pushnotifier.de
func (pn *Pushnotifier) Refresh() error {
	return pn.loginOrRefresh(pn.endpoints.refresh)
}

func (pn *Pushnotifier) loginOrRefresh(endpoint *endpoint) error {
	var (
		statusCode int
		resp       []byte
		err        error
	)

	if endpoint == pn.endpoints.login {
		statusCode, resp, err = pn.apiRequest(endpoint, pn.credentials)
	} else if endpoint == pn.endpoints.refresh {
		statusCode, resp, err = pn.apiRequest(pn.endpoints.refresh, nil)
	}

	if err != nil {
		return err
	}

	if statusCode >= 500 {
		return ErrPushnotifierServerError
	} else if statusCode != 200 {
		return ErrInvalidCredentials
	}

	appToken := &appToken{}

	err = json.Unmarshal(resp, appToken)
	if err != nil {
		return err
	}

	if appToken.Token == "" || appToken.ExpiresAt < 1 {
		return ErrMissingAppToken
	}

	pn.appToken = appToken

	return nil
}

// Devices retrieves all devices from pushnotifier.de
func (pn *Pushnotifier) Devices() ([]*Device, error) {
	_, resp, err := pn.apiRequest(pn.endpoints.devices, nil)
	if err != nil {
		return nil, err
	}

	var devices []*Device

	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return nil, err
	}

	pn.devices = devices

	return pn.GetDevices(), nil
}

// Text sends a text notification to all given devices
// If devices is nil the notification will be send to all devices
func (pn *Pushnotifier) Text(devices []*Device, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return ErrNotificationContentMissing
	}

	notification := &pnNotification{
		Content: content,
	}

	return pn.notification(pn.endpoints.sendText, devices, notification)
}

// URL sends an URL notification to all given devices
// If devices is nil the notification will be send to all devices
func (pn *Pushnotifier) URL(devices []*Device, url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return ErrNotificationURLMissing
	}

	notification := &pnNotification{
		URL: url,
	}

	return pn.notification(pn.endpoints.sendURL, devices, notification)
}

// Notification sends a text and an URL notification to all given devices
// If devices is nil the notification will be send to all devices
func (pn *Pushnotifier) Notification(devices []*Device, content, url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return ErrNotificationURLMissing
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return ErrNotificationContentMissing
	}

	notification := &pnNotification{
		URL:     url,
		Content: content,
	}

	return pn.notification(pn.endpoints.sendURL, devices, notification)
}

// Image sends an image notification to all given devices
// If devices is nil the notification will be send to all devices
func (pn *Pushnotifier) Image(devices []*Device, image []byte, imageName string) error {
	return ErrNotificationNotSupported

	// size := len(image)
	// if size < 1 {
	// 	return ErrNotificationImageMissing
	// } else if size > 5242879 {
	// 	return ErrNotificationImageTooBig
	// }

	// imageName = strings.TrimSpace(imageName)
	// if imageName == "" {
	// 	return ErrNotificationImageNameMissing
	// }

	// notification := &pnNotification{
	// 	Content:  base64.StdEncoding.EncodeToString(image),
	// 	FileName: imageName,
	// }

	// return pn.notification(pn.endpoints.sendImage, devices, notification)
}

func (pn *Pushnotifier) notification(endpoint *endpoint, devices []*Device, notification *pnNotification) error {
	deviceIDs := []string{}

	if devices == nil {
		devices = pn.GetDevices()
	}

	for _, device := range devices {
		deviceIDs = append(deviceIDs, device.ID)
	}

	notification.Devices = deviceIDs

	statusCode, resp, err := pn.apiRequest(endpoint, notification)

	if err != nil {
		return err
	}

	if statusCode >= 500 {
		return ErrPushnotifierServerError
		// } else if statusCode == 413 {
		// 	return ErrNotificationImageTooBig
	} else if statusCode != 200 {
		return ErrDeviceNotFound
	}

	response := &pnNotificationResponse{}
	err = json.Unmarshal(resp, response)
	if err != nil {
		return err
	}

	if len(response.Error) > 0 || len(response.Success) < len(devices) {
		return ErrNotAllDevicesReached
	}

	return nil
}

// GetDevices returns all cached devices which were retrieved in an earlier request
func (pn *Pushnotifier) GetDevices() []*Device {
	return pn.devices
}

func (pn *Pushnotifier) apiRequest(endpoint *endpoint, body interface{}) (int, []byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	url := "https://api.pushnotifier.de/v2/" + endpoint.uri

	var request *http.Request

	switch endpoint.method {
	case "GET":
		request, err = http.NewRequest(http.MethodGet, url, nil)

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

	resp, err := httpClient.Do(request)
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
