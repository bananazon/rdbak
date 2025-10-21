package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gdanko/rdbak/pkg/cookie_jar"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/gdanko/rdbak/pkg/util"
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

func (ac *APIClient) Login(email, pass string) error {
	payload := map[string]any{"email": email, "password": pass}
	payloadStr, _ := json.Marshal(payload)

	loginUrl := url.URL{
		Scheme: "https",
		Host:   apiBase,
		Path:   fmt.Sprintf("%s/auth/email/login", apiVersion),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", loginUrl.String(), bytes.NewBuffer(payloadStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status at login: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginRes data.ListRaindropsResult
	err = json.Unmarshal(body, &loginRes)
	if err != nil {
		return err
	}

	if !loginRes.Result {
		return fmt.Errorf("Login returned false: %s", loginRes.ErrorMessage)
	}

	return nil
}

func (ac *APIClient) ListRaindrops(page int) (listRaindropsResult data.ListRaindropsResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	listUrl := url.URL{
		Scheme:   "https",
		Host:     apiBase,
		Path:     fmt.Sprintf("%s/raindrops/0", apiVersion),
		RawQuery: fmt.Sprintf("sort=-lastUpdate&perpage=%d&page=%d&version=2", PageSize, page),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", listUrl.String(), nil)
	if err != nil {
		return listRaindropsResult, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return listRaindropsResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return listRaindropsResult, fmt.Errorf("bad status at list bookmarks: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listRaindropsResult, err
	}

	err = json.Unmarshal(body, &listRaindropsResult)
	if err != nil {
		return listRaindropsResult, err
	}

	if !listRaindropsResult.Result {
		return listRaindropsResult, fmt.Errorf("list bookmarks returned false: %s", listRaindropsResult.ErrorMessage)
	}

	return listRaindropsResult, nil
}

func (ac *APIClient) ListCollections() (listCollectionsResult data.ListCollectionsResult, err error) {
	// const collectionsUrl = "https://api.raindrop.io/v1/collections"
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	listUrl := url.URL{
		Scheme: "https",
		Host:   apiBase,
		Path:   fmt.Sprintf("%s/collections", apiVersion),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", listUrl.String(), nil)
	if err != nil {
		return listCollectionsResult, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return listCollectionsResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return listCollectionsResult, fmt.Errorf("bad status at list bookmarks: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listCollectionsResult, err
	}

	err = json.Unmarshal(body, &listCollectionsResult)
	if err != nil {
		return listCollectionsResult, err
	}

	if !listCollectionsResult.Result {
		return listCollectionsResult, fmt.Errorf("list bookmarks returned false: %s", listCollectionsResult.ErrorMessage)
	}

	return listCollectionsResult, nil
}

func LimitLength(filename string, maxLen int) string {
	filenameLen := len(filename)
	if filenameLen <= maxLen {
		return filename
	}
	dotix := strings.LastIndex(filename, ".")
	if dotix == -1 {
		return filename[:maxLen]
	}
	extLen := filenameLen - dotix
	if extLen >= maxLen {
		return filename[:maxLen]
	}
	res := filename[:maxLen-extLen] + filename[dotix:]
	return res
}

func safeDeleteFile(filename string) (err error) {
	if util.PathExists(filename) {
		return nil
	}

	if err = os.Remove(filename); err != nil {
		return fmt.Errorf("tried to delete %s, got error: %s", filename, err.Error())
	}

	return nil
}

func (ac *APIClient) getFileName(title string, id uint64, resp *http.Response) string {
	// Baseline: file name is ID
	filename := fmt.Sprintf("%d", id)

	// Download file name is expected in a header
	if cdp := resp.Header.Get("Content-Disposition"); cdp != "" {
		groups := ac.ReDownloadName.FindStringSubmatch(cdp)
		if groups != nil {
			name := LimitLength(groups[1], maxFileNameLen)
			filename += "-" + name
		}
	}

	// Add extension based on mime type, if present
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		if strings.HasPrefix(ct, "application/pdf") {
			filename = strings.TrimSuffix(filename, ".html")
			filename += ".pdf"
		} else if strings.HasPrefix(ct, "text/html") {
			if !strings.HasSuffix(filename, ".html") {
				filename += ".html"
			}
		}
	}

	// If filename is NOT in the header, use title
	if filename == fmt.Sprintf("%d", id) {
		filename = title
	}

	// Whee
	return filename
}

func (ac *APIClient) DownloadFileIfMissing(title string, id uint64, exportDir string) (bool, error) {
	// We need to create a directory for each ID because if two IDs share the same filename,
	// bad things can happen.
	ac.Logger.Infof("Downloading bookmark for ID %d", id)

	etc := util.NewExtensibleTimeoutContext(timeoutSec)
	defer etc.Cancel()

	downloadUrl := url.URL{
		Scheme:   "https",
		Host:     apiBase,
		Path:     fmt.Sprintf("%s/raindrop/%d/cache", apiVersion, id),
		RawQuery: "download",
	}
	url := downloadUrl.String()

	req, _ := http.NewRequestWithContext(etc.Context(), "GET", url, nil)

	resp, err := ac.Client.Do(req)
	if err != nil {
		ac.Logger.Errorf("Error creating client for %s: %s", url, err.Error())
		return false, err
	}
	defer resp.Body.Close()

	// If we don't get a 200 we don't panic. Maybe problem is transient and download
	// will work next time
	if resp.StatusCode != http.StatusOK {
		ac.Logger.Errorf("Got status %d trying to download %s", resp.StatusCode, url)
		return false, err
	}

	filename := ac.getFileName(title, id, resp)
	targetDir := filepath.Join(exportDir, fmt.Sprintf("%d", id))
	filename = filepath.Join(targetDir, filename)

	if util.PathExists(filename) && util.FileSize(filename) != 0 {
		ac.Logger.Infof("File exists; skipping: %s", filename)
		return false, nil
	}

	err = util.VerifyDirectory(targetDir)
	if err != nil {
		return false, err
	}

	outf, err := os.Create(filename)
	if err != nil {
		return false, err
	}
	defer outf.Close()
	ac.Logger.Infof("Saving %s", filename)

	buf := make([]byte, 32*1024)
	savedBytes := 0
	for {
		n, err := resp.Body.Read(buf)
		if (err == nil || err == io.EOF) && n > 0 {
			if _, wrerr := outf.Write(buf[:n]); wrerr != nil {
				ac.Logger.Errorf("Error writing to file: %s", wrerr)
				outf.Close()
				deleteErr := safeDeleteFile(filename)
				if deleteErr != nil {
					ac.Logger.Error(deleteErr)
				}
				return false, err
			}
			savedBytes += n
			etc.Extend()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			ac.Logger.Errorf("Error reading content from %s: %s", url, err)
			outf.Close()
			deleteErr := safeDeleteFile(filename)
			if deleteErr != nil {
				ac.Logger.Error(deleteErr)
			}
			return false, err
		}
	}
	ac.Logger.Infof("Download finished: %d bytes", savedBytes)

	return true, nil
}
