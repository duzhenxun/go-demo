package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-demo/01/jwt"
	"strings"
)

func main() {
	//jwt.io 密钥 123456 注意payload的字符串要与官网一样
	jwtIo:="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.mUC1eX0GwtVeRxXr5a2LvCHrSLMu1YfQ5eXwV5RH5M8"
	payload:=map[string]interface{}{"iat": 1516239022, "name": "John Doe", "sub": "1234567890"}
	p, _ := json.Marshal(payload)
	fmt.Println(jwtEncode(string(p),"123456"))
	t, _ := jwt.Encode(payload, []byte("123456"), "HS256")
	fmt.Println(string(t))
	fmt.Println(jwtIo)
}

func jwtEncode(payload string, secret string) string {
	header, _ := json.Marshal(map[string]interface{}{
		"typ": "JWT",
		"alg": "HS256",
	})
	segments := [3]string{}
	segments[0] = base64url_encode(string(header))
	segments[1] = base64url_encode(payload)

	sha := hmac.New(sha256.New, []byte(secret))
	s:=strings.Join(segments[:2], ".")

	sha.Write([]byte(s))
	segments[2] = base64url_encode(string(sha.Sum(nil)))

	return strings.Join(segments[:], ".")

}

func base64url_encode(b string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(b))
	var equalIndex = strings.Index(encoded, "=")
	if equalIndex > -1 {
		encoded = encoded[:equalIndex]
	}
	return encoded
}
