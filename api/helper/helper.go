package helper

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http:/" + url
	}
	return url
}
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN"){
		return false
	}
	finalURL:=strings.Replace(url,"http:/","",1)
	finalURL=strings.Replace(finalURL,"https:/","",1)
	finalURL=strings.Replace(finalURL,"www.","",1)
	finalURL=strings.Split(finalURL,"/")[0]


	if finalURL == os.Getenv("DOMAIN"){
		return false
	}
	
	return true
}