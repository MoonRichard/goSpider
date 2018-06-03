package main

import (
	"regexp"
	"fmt"
	"net/http"
	"io/ioutil"
	"mahonia"
	"io"
	"os"
)

//<tr class='villagetr'><td>110107002011</td><td>111</td><td>老山东里北社区居委会</td></tr>
var(
	baseURL = `http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/`
	//td for index
	regTd4Index = `<td><a href='\w+\.html'>([\s\S]+?)</a></td>`
	//td标签forNotIndex
	regTr = `<tr class='((citytr)|(countytr)|(towntr)|(villagetr))'>([\s\S]+?)</tr>`
	//中文
	regChar = `[^\x00-\xff]+`
	regCharTitle = `([^\x00-\xff]+代码)|名称`
	//数字
	regNum = `[0-9]{12}`
	//href
	regHref = `\d+(/\d+)?\.html`
	//提取无连接终止页
	regHtmlNoLink = `<tr class='\w+'><td>\d+([\s\S]*)</table>`
	f,_ = os.OpenFile("./data/info.txt",os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
	count = 0
)

func main() {
	parseIndex(getHtml(baseURL))
	f.Close()
}

func parseIndex(html string){
	prvcLabelList := getProvinceList(html)
	for _, p :=range prvcLabelList {
		prvcName := getName(p)
		count++
		fmt.Printf("%s(%d)",prvcName,count)
		io.WriteString(f,prvcName)
		sub := getUrl(p)
		if sub !="" {
			parseCitys(mixUrl(baseURL,sub),"  ")
		}
	}
}

func parseCitys(url string,prefix string){
	html := getHtml(url)
	labelList := getCitysList(html)
	for _, l :=range labelList {
		id,name := getIdAndName(l)
		count++
		fmt.Printf("%s%s %s(%d)\n",prefix,id,name,count)
		str := prefix+id+" "+name+"\n"
		io.WriteString(f,str)
		sub := getUrl(l)
		if sub !="" {
			parseCitys(mixUrl(url,sub),prefix+"  ")
		}
	}
}

func getHtml(url string)  (utf8 string){
	resp,_ := http.Get(url)
	body,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//转码了解一下烦死了
	utf8 = mahonia.NewDecoder("gbk").ConvertString(string(body))
	return
}

func mixUrl(urlFather string,sub string) (url string) {
	regurl := regexp.MustCompile(`\w+\.html`)
	regres := regurl.FindAllString(urlFather,1)
	if regres != nil {
		n := len(urlFather)-len(regres[0])
		rs := []rune(urlFather)
		url = string(rs[0:n])+sub
	}else {
		url = urlFather+sub
	}
	return
}

func getCitysList(html string) (tr []string){
	regTr := regexp.MustCompile(regTr)
	tr = regTr.FindAllString(html,-1)
	return
}

func getProvinceList(html string) (td []string) {
	regTd := regexp.MustCompile(regTd4Index)
	td = regTd.FindAllString(html,-1)
	return
}




func getIdAndName(tr string)(id string,name string){
	regChar := regexp.MustCompile(regChar)
	regNum := regexp.MustCompile(regNum)
	id = regNum.FindString(tr)
	name = regChar.FindString(tr)
	return
}

func getName(tr string) (name string) {
	regChar := regexp.MustCompile(regChar)
	name = regChar.FindString(tr)
	return
	return
}

func getUrl(tr string)(url string){
	regHref := regexp.MustCompile(regHref)
	url = regHref.FindString(tr)
	return
}