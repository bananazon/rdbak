package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bananazon/raindrop/pkg/util"
)

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
			name := limitLength(groups[1], maxFileNameLen)
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

func limitLength(filename string, maxLen int) string {
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
