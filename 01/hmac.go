package main

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
)



func main() {
	src := []byte("hello")
	secret := []byte("123456")

	hasher := hmac.New(md5.New, secret)
	hasher.Write(src)
	mac := hasher.Sum(nil)
	fmt.Printf("%x", mac)
	//68656c6c6fac28d602c767424d0c809edebf73828bed5ce99ce1556f4df8e223faeec60edd
}
