package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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
// 通过正则分离域名，提取二级域名，主域名，端口号
func (*DuduUrl) GetDomainByRegexp(url string, domainList []string) (string, string, string) {
	var hostName string
	var hostPort string
	var subDomain string

	if len(domainList) < 1 {
		return "","",""
	}
	// 正则提取主域名
	regexpStr := "((" + strings.Join(domainList,"|") + ")(/|:|$))"
	domainRegexp,err := regexp.Compile(regexpStr)
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	domainRes := domainRegexp.FindStringSubmatch(url)
	if len(domainRes) != 4 {
		return "", "", ""
	} else {
		hostName = domainRes[2]
		domainStartNum := domainRegexp.FindStringIndex(url)
		if len(domainStartNum) > 0 && domainStartNum[0] >= 1{
			// 提取子域名
			subDomain = url[0:domainStartNum[0]-1]
		}
	}

	portRegexp,_ := regexp.Compile(":([0-9]+)")
	portArr := portRegexp.FindStringSubmatch(url)
	if len(portArr) >= 2 {
		hostPort = portArr[1]
	}

	return subDomain, hostName, hostPort
}
//百度ping提交
func (*DuduUrl) GetBaiDu(domain string, urls []string) bool {
	for _,v := range urls {
		apiUrl := "http://api.share.baidu.com/s.gif?r=" + domain + "&l=" + v
		req,_ := http.NewRequest("GET",apiUrl,nil)
		client := http.Client{}
		_, _ = client.Do(req)
	}

	return true
}

//百度主动提交
type BaiDuData struct {
	Remain			int 		`json:"remain"`
	Success			int 		`json:"success"`
}
func (*DuduUrl) PostBaiDu(dataType int, domain, token string, urls []string) (int,int) {
	var tmpRes BaiDuData
	//类型，0|百度PC，1|百度mip，2|百度熊掌号天级推送，3|百度熊掌号周级推送
	apiUrl := ""
	switch dataType {
	case 0:
		//百度PC推送
		//{"remain":99999,"success":1}
		apiUrl = "http://data.zz.baidu.com/urls?site=" + domain + "&token=" + token
		break
	case 1:
		//百度Mip推送
		//{"remain":99,"success":1,"success_mip":1,"remain_mip":99}
		apiUrl = "http://data.zz.baidu.com/urls?site=" + domain + "&token=" + token + "&type=mip"
		break
	case 2:
		//百度熊掌号推送
		//{"remain":10,"success":0,"not_same_site":["http://www.zehoque.com/2.html"],"success_realtime":0,"remain_realtime":10}
		apiUrl = "http://data.zz.baidu.com/urls?appid=" + domain + "&token=" + token + "&type=realtime"
		break
	case 3:
		//百度熊掌号推送
		//{"remain":5000000,"success":0,"not_same_site":["http://www.zehoque.com/2.html"],"success_batch":0,"remain_batch":5000000}
		apiUrl = "http://data.zz.baidu.com/urls?appid=" + domain + "&token=" + token + "&type=batch"
		break
	default:
		apiUrl = "http://data.zz.baidu.com/urls?site=" + domain + "&token=" + token
		break
	}
	tmpUrl := ""
	for k,v := range urls {
		tmpUrl = v
		if k != len(urls) - 1 {
			tmpUrl = tmpUrl + "\n"
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(tmpUrl))
	if err != nil {
		return 0, 0
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0
	}

	err = json.Unmarshal(body, &tmpRes)
	if err != nil {
		return 0, 0
	}

	return tmpRes.Remain, tmpRes.Success
}

//百度主动提交
type SmData struct {
	ReturnCode		int 		`json:"returnCode"`
}
func (*DuduUrl) PostSm(userName, domain, token string, urls []string) int {
	var tmpRes SmData
	apiUrl := "http://data.zhanzhang.sm.cn/push?site=" + domain + "&user_name=" + userName + "&resource_name=mip_add&token=" + token

	tmpUrl := ""
	for k,v := range urls {
		tmpUrl = v
		if k != len(urls) - 1 {
			tmpUrl = tmpUrl + "\n"
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(tmpUrl))
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	err = json.Unmarshal(body, &tmpRes)
	if err != nil {
		return 0
	}

	return tmpRes.ReturnCode
}