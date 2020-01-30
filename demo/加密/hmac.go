package 加密

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func main() {

	token:=hmacSha256("duzhenxun", "123456")

	fmt.Println(token)
}

//php > echo base64_encode(hash_hmac('sha256','hello','duzhenxun'));
//ZWZkZWM2OGUyNDBiOTg5MDY3ZWQ2ZjAwOGZmOGNhNjQ4ZTA1NTYzYTU2NmZiZTJhMGQ2MzM0MGNjNjM2MGJmNg==
func hmacSha256(src string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(src))
	shaStr:= fmt.Sprintf("%x",h.Sum(nil))
	//shaStr:=hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(shaStr))
}
