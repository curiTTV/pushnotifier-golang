package main

import "errors"

var (
	// ErrMissingPackageName is returned if the package name is empty
	ErrMissingPackageName = errors.New("pushnotifier-golang: package name is missing")

	// ErrMissingAPIToken is returned if the api token is empty
	ErrMissingAPIToken = errors.New("pushnotifier-golang: api token is missing")

	// ErrMissingCredentials is returned if username or password is empty
	ErrMissingCredentials = errors.New("pushnotifier-golang: credentials are missing")

	// ErrInvalidCredentials is returned if the given credentials are invalid for pushnotifier.de
	ErrInvalidCredentials = errors.New("pushnotifier-golang: invalid credentials provided")

	// ErrPushnotifierServerError is returned if a server error >500 happened on any request
	ErrPushnotifierServerError = errors.New("pushnotifier-golang: server error happened on pushnotifier.de")

	// ErrInvalidHTTPMethod is returned if a http method is used which is not supported by pushnotifier.de
	ErrInvalidHTTPMethod = errors.New("pushnotifier-golang: invalid http method provided")
)

type (
	// Pushnotifier holds every necessary method to interact with pushnotifier.de
	Pushnotifier struct {
		packageName string
		apiToken    string
		credentials *credentials
		appToken    *appToken

		endpoints map[string]*endpoint

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
)
