package api

import (
	"net/http"
	"regexp"

	"github.com/gdanko/rdbak/pkg/cookie_jar"
	"github.com/sirupsen/logrus"
)

const apiBase = "api.raindrop.io"
const apiVersion = "v1"
const PageSize = 40
const maxFileNameLen = 128
const timeoutSec = 60

// const collectionsChildrenUrl = "https://api.raindrop.io/v1/collections/childrens"
// const collectionUrl = "https://api.raindrop.io/v1/collection/{id}"

type APIClient struct {
	Jar            *cookie_jar.CookieJar
	Logger         *logrus.Logger
	Client         *http.Client
	ReDownloadName *regexp.Regexp
}

func NewApiClient(logger *logrus.Logger) *APIClient {
	ac := APIClient{
		Logger: logger,
	}
	ac.Jar = cookie_jar.NewJar()
	ac.Client = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           ac.Jar,
		Timeout:       0,
	}
	ac.ReDownloadName = regexp.MustCompile("attachment; filename=\"(.+)\"")
	return &ac
}
