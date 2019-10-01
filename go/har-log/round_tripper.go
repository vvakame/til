package harlog

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
)

var _ http.RoundTripper = (*Transport)(nil)

// Transport is collecting http request/response log by HAR format.
type Transport struct {
	// next Transport. if nil, use http.DefaultTransport.
	Transport http.RoundTripper
	// unusual (not network oriented) error occurred, handle error by this function.
	// if nil, emit error log by log package, and ignore it.
	UnusualError func(err error) error

	har   *HARContainer
	mutex sync.Mutex
}

func (h *Transport) init() {
	if h.har != nil {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.har != nil {
		return
	}

	h.har = &HARContainer{
		Log: &Log{
			Version: "1.2",
			Creator: &Creator{
				Name:    "github.com/vvakame/til/go/har-log",
				Version: "0.0.1",
			},
		},
	}
}

// HAR returns HAR format log data.
func (h *Transport) HAR() *HARContainer {
	h.init()
	return h.har
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (h *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	h.init()

	baseRoundTripper := h.Transport
	if baseRoundTripper == nil {
		baseRoundTripper = http.DefaultTransport
	}

	entry := &Entry{}
	defer func() {
		h.mutex.Lock()
		h.har.Log.Entries = append(h.har.Log.Entries, entry)
		h.mutex.Unlock()
	}()

	err := h.preRoundTrip(r, entry)
	if err != nil {
		if h.UnusualError != nil {
			err = h.UnusualError(err)
		} else {
			log.Println(err)
			err = nil
		}
		if err != nil {
			return nil, err
		}
	}

	trace, ct, finish := newClientTracer()
	r = r.WithContext(httptrace.WithClientTrace(r.Context(), ct))
	defer func() {
		entry.StartedDateTime = Time(trace.startAt)
		entry.Time = Duration(trace.endAt.Sub(trace.startAt))
		entry.Timings = &Timings{
			Blocked: Duration(trace.startAt.Sub(trace.connStart)),
			DNS:     -1,
			Connect: -1,
			Send:    Duration(trace.writeRequest.Sub(trace.connObtained)),
			Wait:    Duration(trace.firstResponseByte.Sub(trace.writeRequest)),
			Receive: Duration(trace.endAt.Sub(trace.firstResponseByte)),
			SSL:     -1,
		}
		if !trace.dnsStart.IsZero() {
			entry.Timings.DNS = Duration(trace.dnsEnd.Sub(trace.dnsStart))
		}
		if !trace.connStart.IsZero() {
			entry.Timings.Connect = Duration(trace.connObtained.Sub(trace.connStart))
		}
		if !trace.tlsHandshakeStart.IsZero() {
			entry.Timings.SSL = Duration(trace.tlsHandshakeEnd.Sub(trace.tlsHandshakeStart))
		}
	}()

	resp, realErr := baseRoundTripper.RoundTrip(r)

	err = h.postRoundTrip(r, resp, entry, finish)
	if err != nil {
		if h.UnusualError != nil {
			err = h.UnusualError(err)
		} else {
			log.Println(err)
			err = nil
		}
		if err != nil {
			return nil, err
		}
	}

	entry.Cache = &Cache{}

	return resp, realErr
}

func (h *Transport) preRoundTrip(r *http.Request, entry *Entry) error {
	bodySize := -1
	var postData *PostData
	if r.Body != nil {
		reqBody, err := r.GetBody()
		if err != nil {
			return err
		}

		reqBodyBytes, err := ioutil.ReadAll(reqBody)
		if err != nil {
			return err
		}
		bodySize = len(reqBodyBytes)

		mimeType := r.Header.Get("Content-Type")
		postData = &PostData{
			MimeType: mimeType,
			Params:   []*Param{},
			Text:     string(reqBodyBytes),
		}

		mediaType, _, err := mime.ParseMediaType(mimeType)
		if err != nil {
			return err
		}

		switch mediaType {
		case "application/x-www-form-urlencoded":
			err := r.ParseForm()
			if err != nil {
				return err
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

			for k, v := range r.PostForm {
				for _, s := range v {
					postData.Params = append(postData.Params, &Param{
						Name:  k,
						Value: s,
					})
				}
			}

		case "multipart/form-data":
			err := r.ParseMultipartForm(10 * 1024 * 1024)
			if err != nil {
				return err
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

			for k, v := range r.MultipartForm.Value {
				for _, s := range v {
					postData.Params = append(postData.Params, &Param{
						Name:  k,
						Value: s,
					})
				}
			}
			for k, v := range r.MultipartForm.File {
				for _, s := range v {
					postData.Params = append(postData.Params, &Param{
						Name:        k,
						FileName:    s.Filename,
						ContentType: s.Header.Get("Content-Type"),
					})
				}
			}
		}
	}

	entry.Request = &Request{
		Method:      r.Method,
		URL:         r.URL.String(),
		HTTPVersion: r.Proto,
		Cookies:     h.toHARCookies(r.Cookies()),
		Headers:     h.toHARNVP(r.Header),
		QueryString: h.toHARNVP(r.URL.Query()),
		PostData:    postData,
		HeadersSize: -1, // TODO
		BodySize:    bodySize,
	}

	return nil
}

func (h *Transport) postRoundTrip(r *http.Request, resp *http.Response, entry *Entry, finish func()) error {
	if resp == nil {
		finish()
		return nil
	}
	respBody := resp.Body
	respBodyBytes, err := ioutil.ReadAll(respBody)
	defer func() {
		_ = respBody.Close()
	}()
	finish() // データ読み終わった瞬間が終わり
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))

	mimeType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return err
	}
	var text string
	var encoding string
	switch {
	case strings.HasPrefix(mediaType, "text/"):
		text = string(respBodyBytes)
	default:
		text = base64.StdEncoding.EncodeToString(respBodyBytes)
		encoding = "base64"
	}

	entry.Response = &Response{
		Status:      resp.StatusCode,
		StatusText:  "",
		HTTPVersion: resp.Proto,
		Cookies:     h.toHARCookies(resp.Cookies()),
		Headers:     h.toHARNVP(resp.Header),
		Content: &Content{
			Size:        resp.ContentLength, // TODO 圧縮されている場合のフォロー
			Compression: 0,
			MimeType:    mimeType,
			Text:        text,
			Encoding:    encoding,
		},
		RedirectURL: resp.Header.Get("Location"),
		HeadersSize: -1,
		BodySize:    resp.ContentLength,
	}

	return nil
}

func (h *Transport) toHARCookies(cookies []*http.Cookie) []*Cookie {
	harCookies := make([]*Cookie, 0, len(cookies))

	for _, cookie := range cookies {
		harCookies = append(harCookies, &Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  Time(cookie.Expires),
			HTTPOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		})
	}

	return harCookies
}

func (h *Transport) toHARNVP(vs map[string][]string) []*NVP {
	nvps := make([]*NVP, 0, len(vs))

	for k, v := range vs {
		for _, s := range v {
			nvps = append(nvps, &NVP{
				Name:  k,
				Value: s,
			})
		}
	}

	return nvps
}
