package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// USER_AGENT http header User-Agent's
var USER_AGENTS = [8]string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A",
	"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5355d Safari/8536.25",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.13+ (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/534.55.3 (KHTML, like Gecko) Version/5.1.3 Safari/534.53.10",
}

// RandomUserAgent generate random User-Agent
func RandomUserAgent() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return USER_AGENTS[r1.Intn(len(USER_AGENTS))]
}

// DefaultHttpClient default http client use http.Client
// timeout 5 Second
// ignore certificate warnings
var DefaultHttpClient = &http.Client{Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore certificate warnings
	},
	CheckRedirect: nil}

var NoneCheckRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

// IsDomainName checks if a string is a presentation-format domain name
// (currently restricted to hostname-compatible "preferred name" LDH labels and
// SRV-like "underscore labels"; see golang.org/issue/12421).
func IsDomainName(s string) bool {
	// The root domain name is valid. See golang.org/issue/45715.
	if s == "." {
		return true
	}

	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return nonNumeric
}

type HttpClient struct {
	http.Client

	header http.Header
}

func (c *HttpClient) PreDo(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if c.header != nil {
		req.Header = c.header
	}

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	return c.PreDo("GET", url, nil)
}

func (c *HttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	c.Headers("Content-Type", contentType)
	return c.PreDo("POST", url, body)
}

func (c *HttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	c.Headers("Content-Type", "application/x-www-form-urlencoded")
	return c.PreDo("Post", url, strings.NewReader(data.Encode()))
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type UPFile struct {
	Path        string
	Field       string
	Filename    string
	ContentType string
}

const DefaultFileContentType = "application/octet-stream"

// UploadFile ref https://gist.github.com/andrewmilson/19185aab2347f6ad29f5
func (c *HttpClient) UploadFile(url string, files []UPFile) (resp *http.Response, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := range files {
		op, err := os.Open(files[i].Path)
		if err != nil {
			return nil, err
		}

		filename := files[i].Filename
		if filename == "" {
			filename = filepath.Base(op.Name())
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(files[i].Field), escapeQuotes(filename)))
		if files[i].ContentType == "" {
			h.Set("Content-Type", DefaultFileContentType)
		} else {
			h.Set("Content-Type", files[i].ContentType)
		}

		part, _ := writer.CreatePart(h)

		io.Copy(part, op)
	}

	c.Headers("Content-Type", writer.FormDataContentType())

	writer.Close()

	return c.PreDo("POST", url, body)
}

func (c *HttpClient) Headers(k string, v string) {
	if c.header == nil {
		c.header = make(http.Header, 0)
	}
	c.header.Set(k, v)
}
