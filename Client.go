package nicehash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Client struct {
	orgID  string
	key    string
	secret string

	DumpWriter io.Writer
}

func NewClient(OrgID string, key string, secret string) *Client {
	return &Client{
		orgID:  OrgID,
		key:    key,
		secret: secret,
	}
}

func (c *Client) dump(reqOrResponse interface{}) {
	if c.DumpWriter == nil {
		return
	}

	switch r := reqOrResponse.(type) {
	case *http.Request:
		d, _ := httputil.DumpRequestOut(r, true)
		_, _ = c.DumpWriter.Write(d)

	case *http.Response:
		d, _ := httputil.DumpResponse(r, true)
		_, _ = c.DumpWriter.Write(d)
	}
}

func (c *Client) signRequest(req *http.Request) {
	mac := hmac.New(sha256.New, []byte(c.secret))
	sep := false

	writeString := func(str string) {
		if sep {
			_, _ = mac.Write([]byte{0})
		}

		// Writing to a HMAC should never fail.
		_, _ = mac.Write([]byte(str))

		sep = true
	}

	writeString(c.key)
	writeString(req.Header.Get("X-Time"))
	writeString(req.Header.Get("X-Nonce"))
	writeString("")
	writeString(req.Header.Get("X-Organization-Id"))
	writeString("")
	writeString(strings.ToUpper(req.Method))
	writeString(req.URL.Path)
	writeString(req.URL.RawQuery)

	if req.Body != nil {
		writeString("")

		// This should never fail for our own requests.
		body, _ := req.GetBody()
		_, _ = io.Copy(mac, body)
	}

	req.Header.Add("X-Auth", c.key+":"+hex.EncodeToString(mac.Sum(nil)))
}

func (c *Client) exchangeJSON(method string, url string, payload interface{}, target interface{}) error {
	body := io.Reader(nil)

	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(encoded)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Add("content-type", "application/json")
	}

	rand := randString(30)

	unix := time.Now().UTC().Unix()
	unixStr := fmt.Sprintf("%d", unix*1000)

	req.Header.Add("accept", "application/json")

	// Headers requires by the Nicehash API.
	req.Header.Add("X-Time", unixStr)
	req.Header.Add("X-Nonce", rand)
	req.Header.Add("X-Organization-Id", c.orgID)
	req.Header.Add("X-Request-Id", rand)

	c.signRequest(req)

	c.dump(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.dump(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Unexpected status from %s: %d", url, resp.StatusCode)
	}

	if target == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
