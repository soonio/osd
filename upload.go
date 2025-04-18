package osd

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
	"strings"
	"time"
)

type Policy struct {
	AccessId  string `json:"accessId"`
	Type      string `json:"type"`
	Host      string `json:"host"`
	Policy    string `json:"policy"`
	Signature string `json:"signature"`
	Expire    int64  `json:"expire"`
	Callback  string `json:"callback"`
	Dir       string `json:"dir"`
}

type Client struct {
	key      string // key
	secret   string // 密钥
	host     string // 域名
	prefix   string // 缓存前缀
	callback string // 回调地址
	duration int64  // 生成的上传策略有效时长
}

func New(key, secret, host string, opt ...Applies) *Client {
	c := &Client{
		key:      key,
		secret:   secret,
		host:     host,
		prefix:   "",
		duration: 120,
	}
	for _, o := range opt {
		o(c)
	}
	return c
}

func (o *Client) Signature(f string) *Policy {
	var dir = strings.ReplaceAll(strings.TrimLeft(o.prefix+f, "/"), "//", "/")

	t := time.Now().Add(time.Second * time.Duration(o.duration))

	policy, _ := json.Marshal(map[string]any{
		"expiration": t.UTC().Format("2006-01-02T15:04:05Z"),
		"conditions": [][]string{{"starts-with", "$key", dir}},
	})
	base64Policy := base64.StdEncoding.EncodeToString(policy)

	signature := base64.StdEncoding.EncodeToString(hashHmac(base64Policy, o.secret))

	return &Policy{
		AccessId:  o.key,
		Host:      o.host,
		Policy:    base64Policy,
		Signature: signature,
		Expire:    t.Unix(),
		Callback:  callback(o.callback),
		Dir:       dir,
	}
}

func callback(url string) string {
	if url != "" {
		var m = map[string]string{
			"callbackUrl":      url,
			"callbackBody":     "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}",
			"callbackBodyType": "application/json",
		}
		bs, _ := json.Marshal(m)
		return base64.StdEncoding.EncodeToString(bs)
	}
	return ""
}

func hashHmac(str string, key string) []byte {
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(key))
	_, _ = io.WriteString(h, str)
	return h.Sum(nil)
}
