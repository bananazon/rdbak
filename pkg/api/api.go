package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gdanko/rdbak/pkg/cookie_jar"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/gdanko/rdbak/util"
	"github.com/sirupsen/logrus"
)

const PageSize = 40
const maxFileNameLen = 128
const timeoutSec = 60
const loginUrl = "https://api.raindrop.io/v1/auth/email/login"
const listUrl = "https://api.raindrop.io/v1/raindrops/0?sort=-lastUpdate&perpage=%v&page=%v&version=2"
const downloadUrl = "https://api.raindrop.io/v1/raindrop/%v/cache?download"
const collsUrl = "https://api.raindrop.io/v1/collections"
const collsChildrenUrl = "https://api.raindrop.io/v1/collections/childrens"

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

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", loginUrl, bytes.NewBuffer(payloadStr))
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

	var loginRes data.ListRes
	err = json.Unmarshal(body, &loginRes)
	if err != nil {
		return err
	}
	if !loginRes.Result {
		return fmt.Errorf("Login returned false: %s", loginRes.ErrorMessage)
	}

	return nil
}

func (ac *APIClient) ListBookmarks(page int) (listResult data.ListRes, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()
	url := fmt.Sprintf(listUrl, PageSize, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return listResult, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return listResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return listResult, fmt.Errorf("bad status at list bookmarks: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return listResult, err
	}

	err = json.Unmarshal(body, &listResult)
	if err != nil {
		return listResult, err
	}

	if !listResult.Result {
		return listResult, fmt.Errorf("list bookmarks returned false: %s", listResult.ErrorMessage)
	}

	return listResult, nil
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
	if _, err = os.Stat(filename); err != nil {
		return nil
	}

	if err = os.Remove(filename); err != nil {
		return fmt.Errorf("tried to delete %s, got error: %s", filename, err.Error())
	}

	return nil
}

func (ac *APIClient) getFileName(id uint64, resp *http.Response) string {

	// Baseline: file name is ID
	filename := fmt.Sprintf("%v", id)

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

	// Whee
	return filename
}

func (ac *APIClient) DownloadFileIfMissing(id uint64, dir string) (bool, error) {
	ac.Logger.Infof("Downloading bookmark %d", id)

	etc := util.NewExtensibleTimeoutContext(timeoutSec)
	defer etc.Cancel()
	url := fmt.Sprintf(downloadUrl, id)
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

	filename := ac.getFileName(id, resp)
	filename = path.Join(dir, filename)

	if stat, err := os.Stat(filename); err == nil && stat.Size() != 0 {
		ac.Logger.Infof("File exists; skipping: %s", filename)
		return false, nil
	}

	outf, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	ac.Logger.Infof("Saving %s\n", filename)

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
	ac.Logger.Infof("Download finished: %d bytes\n", savedBytes)

	return true, nil
}
