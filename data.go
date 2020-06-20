package main

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
	ErrPushnotifierServerError = errors.New("pushnotifier-golang: server error happened on pushnotifier.de")

	// ErrInvalidHTTPMethod is returned if a http method is used which is not supported by pushnotifier.de
	ErrInvalidHTTPMethod = errors.New("pushnotifier-golang: invalid http method provided")

	// ErrNotificationContentMissing is returned if no content was provided
	ErrNotificationContentMissing = errors.New("pushnotifier-golang: empty content provided")

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
		Devices []string `json:"devices"`
		Content string   `json:"content,omitempty"`
		URL     string   `json:"url,omitempty"`
	}

	pnNotificationResponse struct {
		Success []string `json:"success"`
		Error   []string `json:"error"`
	}
)
