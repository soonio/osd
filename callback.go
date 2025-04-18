package osd

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Verify(r *http.Request) bool {
	bytePublicKey, err := getPublicKey(r)
	if err != nil {
		return false
	}
	byteAuthorization, err := getAuthorization(r)
	if err != nil {
		return false
	}
	byteMD5, err := getMD5FromNewAuthString(r)
	if err != nil {
		return false
	}
	return verifySignature(bytePublicKey, byteMD5, byteAuthorization)
}

func verifySignature(bytePublicKey []byte, byteMd5 []byte, authorization []byte) bool {
	pubBlock, _ := pem.Decode(bytePublicKey)
	if pubBlock == nil {
		return false
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if (pubInterface == nil) || (err != nil) {
		return false
	}
	pub := pubInterface.(*rsa.PublicKey)

	err = rsa.VerifyPKCS1v15(pub, crypto.MD5, byteMd5, authorization)
	return err == nil
}

func getPublicKey(r *http.Request) ([]byte, error) {
	var bs []byte

	header := r.Header.Get("x-oss-pub-key-url")
	if header == "" {
		return bs, errors.New("no x-storage-pub-key-url field in Request header ")
	}
	url, _ := base64.StdEncoding.DecodeString(header)

	rps, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}
	bs, err = io.ReadAll(rps.Body)
	if err != nil {
		return bs, err
	}
	defer func(r *http.Response) { _ = r.Body.Close() }(rps)
	return bs, nil
}

func getAuthorization(r *http.Request) ([]byte, error) {
	var bs []byte
	header := r.Header.Get("authorization")
	if header == "" {
		return bs, errors.New("no authorization field in Request header")
	}
	bs, _ = base64.StdEncoding.DecodeString(header)
	return bs, nil
}

func getMD5FromNewAuthString(r *http.Request) ([]byte, error) {
	var bs []byte
	bodyContent, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return bs, err
	}
	strCallbackBody := string(bodyContent)

	strAuth := ""
	if r.URL.RawQuery == "" {
		strAuth = fmt.Sprintf("%s\n%s", r.URL.Path, strCallbackBody)
	} else {
		strAuth = fmt.Sprintf("%s?%s\n%s", r.URL.Path, r.URL.RawQuery, strCallbackBody)
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strAuth))
	bs = md5Ctx.Sum(nil)

	return bs, nil
}
