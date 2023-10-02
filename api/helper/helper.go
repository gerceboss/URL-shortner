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
	newURL:=strings.Replace(url,"http:/","",1)
	newURL1:=strings.Replace(newURL,"https:/","",1)
	newURL2:=strings.Split(newURL1,"/")[0]
	//newURL:=strings.Replace(newURL,"www.","",1)

	
	if newURL2 == os.Getenv("DOMAIN"){
		return false
	}
	
	return true
}