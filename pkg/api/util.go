package api

import "net/url"

func MapToQueryString(m map[string]string) string {
	v := url.Values{}
	for k, val := range m {
		v.Set(k, val)
	}
	return v.Encode()
}
