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

var(
	baseURL = `http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/`
	//td for index
	regTd4Index = `<td><a href='\w+\.html'>([\s\S]+?)</a></td>`
	//td标签forNotIndex
	regTr = `<tr class='((citytr)|(countytr)|(towntr)|(villagetr))'>([\s\S]+?)</tr>`
	//中文
	regChar = `[^\x00-\xff]+`
	//数字
	regNum = `[0-9]{12}`
	//href
	regHref = `\d+(/\d+)?\.html`
	ch  = make(chan string)
)

func main() {
	parseIndex(getHtml(baseURL))
	for i:=0;i<31;i++{
		fmt.Printf("%s完成",<-ch)
	}
}

func test(){
	url:=`http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/51.html`
	prefex := ""
	fileName := "./data/sc.txt"
	provinceName := "四川"
	parseCitys(url,prefex,fileName,provinceName)
}

func parseIndex(html string){
	prvcLabelList := getProvinceList(html)
	for _, p :=range prvcLabelList {
		prvcName := getName(p)
		fileName := "./data/"+prvcName+".txt"
		sub := getUrl(p)
		go parseCitys(mixUrl(baseURL,sub),"",fileName,prvcName)
	}
}

func parseCitys(url string,prefix string,fileName string,provinceName string){
	html := getHtml(url)
	labelList := getCitysList(html)
	f,_ := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_RDWR,0)
	for _, l :=range labelList {
		id,name := getIdAndName(l)
		str := prefix+id+" "+name+"\n"
		fmt.Print(str)
		io.WriteString(f,str)
		sub := getUrl(l)
		if sub !="" {
			parseCitysHelp(mixUrl(url,sub),prefix+"  ",f)
		}

	}
	f.Close()
	ch <- provinceName
}

func parseCitysHelp(url string,prefix string,f *os.File){
	html := getHtml(url)
	labelList := getCitysList(html)
	for _, l :=range labelList {
		id,name := getIdAndName(l)
		str := prefix+id+" "+name+"\n"
		fmt.Print(str)
		io.WriteString(f,str)
		sub := getUrl(l)
		if sub !="" {
			parseCitysHelp(mixUrl(url,sub),prefix+"  ",f)
		}
	}
}

func getHtml(url string) (content string){
	resp,_ := http.Get(url)
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	//转码了解一下烦死了
	content = mahonia.NewDecoder("gbk").ConvertString(string(body))
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