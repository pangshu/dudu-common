package common

import (
	"net"
	"regexp"
	"strings"
)
// 分离域名，提取二级域名，主域名，端口号
func (*DuduUrl) GetDomain(url string) (string,string,string) {
	var hostName string
	var hostPort string
	var subDomain string

	isPort := strings.Contains(url, ":")
	if isPort {
		h := strings.Split(url, ":")
		if len(h) >= 2 {
			hostName,hostPort,_ = net.SplitHostPort(url)
		}
	} else {
		hostName = url
		hostPort = ""
	}

	domainRules,_ := regexp.Compile("[\\w][\\w-]*\\.(?:com\\.cn|com|cn|co|net\\.cn|net|org|gov|cc|biz|info|me|xyz|im|io|name|tw|mobi|asia|hk|areo|ca|us|fr|se|ie|tv|ws)(\\/|$)")
	domains := domainRules.FindString(hostName)

	if len(domains) > 0 {
		subDomain = strings.Replace(hostName, domains, "", -1)
		if len(subDomain) > 0 {
			if subDomain[len(subDomain)-1:] == "." {
				subDomain = subDomain[0:len(subDomain)-1]
			}
		}
	}
	return subDomain, domains, hostPort
}