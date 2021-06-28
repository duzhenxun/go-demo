package code

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	//jwt.io 密钥 123456 注意payload的字符串要与官网一样
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZ2UiOjIwLCJuYW1lIjoiZHV6aGVueHVuIn0.Ddh8knTN5SgLFuk4_04sijv8i906xtzGpjkhJEaaBcA
	payload := map[string]interface{}{"age": 20, "name": "duzhenxun"}
	secret := "123456"
	p, _ := json.Marshal(payload)
	token := jwtEncode(string(p), secret)
	fmt.Println(token)
}

func jwtEncode(payload string, secret string) string {
	header := `{"alg":"HS256","typ":"JWT"}`
	segments := [3]string{}
	segments[0] = base64url_encode(string(header))
	segments[1] = base64url_encode(payload)

	sha := hmac.New(sha256.New, []byte(secret))
	s := strings.Join(segments[:2], ".")

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
