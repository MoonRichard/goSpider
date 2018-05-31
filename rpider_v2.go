package main

import (
	"regexp"
	"fmt"
	"io/ioutil"
	"net/http"
	"mahonia"
	"strconv"
)
var(
	//起始页面
	baseURL = `http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/`
	regHref = `\d+(/\d+)?\.html`
	regList = `<a href=([\s\S]*)</table>`
	//中文
	regChar = `[^\x00-\xff]+`
	//数字
	regNum = `>[0-9]+<`
	count = 0

)

func main(){
	initUrl()
}


func initUrl(){
	resp,_ := http.Get(baseURL)
	body,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//转码了解一下烦死了
	utf8 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	tableStr := getList(utf8)
	//链接list
	provinceList := getName(tableStr)
	urlList := getUrlList4Index(tableStr)
	//打印
	if len(provinceList)==len(urlList) {
		i := 0
		for ;i<len(provinceList) ; i++ {
			fmt.Println(provinceList[i])
			count ++
			next_url:= mixUrl(baseURL,urlList[i])
			parsePage(next_url,"  ")
		}
	}else {
		fmt.Println("解析首页错误")
	}
}

func parsePage(url string,prefix string) {
	resp,_ := http.Get(url)
	body,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//转码了解一下烦死了
	utf8 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	table := getList(utf8)
	idList , nameList := getIdAndName(table)
	urlList := getUrlList(table)
	i := 0
	for ; i<len(idList);i++  {
		count++
		fmt.Print(prefix+idList[i]+" "+nameList[i])
		fmt.Println(""+strconv.Itoa(count))
		next_url := mixUrl(url,urlList[i])
		parsePage(next_url,prefix+"  ")
	}
	return
}

func printIdAndName(id []string,name []string,outPrefix string){
	i := 1
	for ; i<len(id) ;i++{
		fmt.Println(outPrefix+id[i]+" "+name[i])
	}
}

//拼接URL
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

//获取list部分的str为getIdAndName服务
func getList(doc string) (rs string){
	regList := regexp.MustCompile(regList)
	rs = regList.FindString(doc)
	return
}
//获取id-name的列表
func getIdAndName(str string)(id []string,name []string){
	regChar := regexp.MustCompile(regChar)
	regNum := regexp.MustCompile(regNum)
	for _,c := range regNum.FindAllString(str,-1){
		c = string(c[1:len(c)-1])
		id = append(id,c)
	}
	for _,c := range regChar.FindAllString(str,-1){
		name = append(name,c)
	}
	return
}

func getUrlList4Index(str string)(url []string){
	regHref := regexp.MustCompile(regHref)
	for _,c := range regHref.FindAllString(str,-1){
		url = append(url,c)
	}
	return
}

func getUrlList(str string)(url []string){
	regHref := regexp.MustCompile(regHref)
	i := 0
	for _,c := range regHref.FindAllString(str,-1){
		if i%2 == 0 {
			url = append(url,c)
		}
		i++
	}
	return
}

//index用
func getName(str string) (name []string) {
	regChar := regexp.MustCompile(regChar)
	for _,c := range regChar.FindAllString(str,-1){
		name = append(name,c)
	}
	return
}
