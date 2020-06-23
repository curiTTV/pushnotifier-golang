package pushnotifier

import (
	"errors"
	"net/http"
)

var (
	// ErrMissingPackageName is returned if the package name is empty
	ErrMissingPackageName = errors.New("pushnotifier-golang: package name is missing")

	// ErrMissingAPIToken is returned if the api token is empty
	ErrMissingAPIToken = errors.New("pushnotifier-golang: api token is missing")

	// ErrMissingCredentials is returned if username or password is empty
	ErrMissingCredentials = errors.New("pushnotifier-golang: credentials are missing")

	// ErrInvalidCredentials is returned if the given credentials are invalid for pushnotifier.de
	ErrInvalidCredentials = errors.New("pushnotifier-golang: invalid credentials provided")

	// ErrMissingAppToken is returned if either on login or refresh the app token is empty
	ErrMissingAppToken = errors.New("pushnotifier-golang: returned app token empty")

	// ErrPushnotifierServerError is returned if a server error >500 happened on any request
	// This can happen if a notification should be send to a registered device which actually does not "exists" any more
	ErrPushnotifierServerError = errors.New("pushnotifier-golang: server error happened on pushnotifier.de")

	// ErrInvalidHTTPMethod is returned if a http method is used which is not supported by pushnotifier.de
	ErrInvalidHTTPMethod = errors.New("pushnotifier-golang: invalid http method provided")

	// ErrNotificationContentMissing is returned if no content was provided
	ErrNotificationContentMissing = errors.New("pushnotifier-golang: content not provided")

	// ErrNotificationURLMissing is returned if no url was provided
	ErrNotificationURLMissing = errors.New("pushnotifier-golang: url not provided")

	// ErrNotificationImageMissing is returned if the content of an image was not provided
	ErrNotificationImageMissing = errors.New("pushnotifier-golang: image content not provided")

	// ErrNotificationImageTooBig is returned if the the size of the image was too big
	ErrNotificationImageTooBig = errors.New("pushnotifier-golang: image file size is to big")

	// ErrNotificationImageNameMissing is returned if no image name was provided
	ErrNotificationImageNameMissing = errors.New("pushnotifier-golang: image name not provided")

	// ErrNotificationNotSupported is returned if a feature is not supported yet
	ErrNotificationNotSupported = errors.New("pushnotifier-golang: feature not supported yet")

	// ErrDeviceNotFound is returned if pushnotifier.de can not find at least one of the given devices
	ErrDeviceNotFound = errors.New("pushnotifier-golang: device could not be found")

	// ErrNotAllDevicesReached is returned if pushnotifier.de could not send the notification to any given device
	ErrNotAllDevicesReached = errors.New("pushnotifier-golang: notification could not send to all given devices")

	httpClient HTTPClient
)

type (
	// HTTPClient is an interface for http.Client
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// Pushnotifier holds every necessary method to interact with pushnotifier.de
	Pushnotifier struct {
		packageName string
		apiToken    string
		credentials *credentials
		appToken    *appToken

		endpoints struct {
			login            *endpoint
			refresh          *endpoint
			devices          *endpoint
			sendText         *endpoint
			sendURL          *endpoint
			sendNotification *endpoint
			sendImage        *endpoint
		}

		devices []*Device
	}

	credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	appToken struct {
		Token     string `json:"app_token"`
		ExpiresAt int    `json:"expires_at"`
	}

	endpoint struct {
		method string
		uri    string
	}

	// Device represents a device registered on pushnotifier.de
	Device struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Model string `json:"model"`
		Image string `json:"image"`
	}

	pnNotification struct {
		Devices  []string `json:"devices"`
		Content  string   `json:"content,omitempty"`
		URL      string   `json:"url,omitempty"`
		Silent   bool     `json:"silent"`
		FileName string   `json:"filename"`
	}

	pnNotificationResponse struct {
		Success []string `json:"success"`
		Error   []string `json:"error"`
	}
)
