package main
//被金门县搞坏的爬虫
import (
	"regexp"
	"fmt"
	"io/ioutil"
	"net/http"
	"mahonia"
	"strconv"
	"os"
	"io"
)
var(
	//起始页面
	baseURL = `http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/`

	regListv2 = `<table class([\s\S]*)</table>`
	regList = `<a href=([\s\S]*)</table>`
	regTd = ``
	//中文
	regChar = `[^\x00-\xff]+`
	regCharTitle = `([^\x00-\xff]+代码)|名称`
	//数字
	regNum = `[0-9]{12}`
	//href
	regHref = `\d+(/\d+)?\.html`
	count = 0
	f,_ = os.OpenFile("./data/info.txt",os.O_CREATE|os.O_RDWR, 0666)
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
			count ++
			str := provinceList[i]+strconv.Itoa(count)
			io.WriteString(f,str+"\n")
			fmt.Printf("%s (%d)",provinceList[i],count)
			fmt.Println()
			next_url:= mixUrl(baseURL,urlList[i])
			parsePage(next_url,"  ")
		}
	}else {
		fmt.Println("解析首页错误")
	}
	f.Close()
}
func parsePage(url string,prefix string) {
	resp,err := http.Get(url)
	if err!=nil {
		fmt.Println(url)
		return
	}
	body,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//转码了解一下烦死了
	utf8 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	table := getListV2(utf8)
	idList , nameList := getIdAndName(table)
	urlList := getUrlList(table)
	isUrlExist := true
	if len(urlList) == 0 {
		isUrlExist = false
	}
	if len(urlList)!= len(idList) {
		tmp := append([]string{},urlList[0:]...)
		urlList=append(urlList[:0],"0.html")
		urlList=append(urlList,tmp...)
	}
	i := 0
	for ; i<len(idList);i++  {
		count++
		str := prefix+idList[i]+" "+nameList[i]+"("+strconv.Itoa(count)+")"
		io.WriteString(f,str+"\n")
		fmt.Println(prefix+idList[i]+" "+nameList[i]+"("+strconv.Itoa(count)+")")
		if isUrlExist {
			nextUrl := mixUrl(url,urlList[i])
			parsePage(nextUrl,prefix+"  ")
		}
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

func getListV2(doc string) (rs string){
	regList := regexp.MustCompile(regListv2)
	rs = regList.FindString(doc)
	return
}

//获取id-name的列表
func getIdAndName(str string)(id []string,name []string){
	regChar := regexp.MustCompile(regChar)
	regNum := regexp.MustCompile(regNum)
	regTitle := regexp.MustCompile(regCharTitle)
	for _,c := range regNum.FindAllString(str,-1){
		c = string(c)
		id = append(id,c)
	}
	for _,c := range regChar.FindAllString(str,-1){
		if regTitle.FindString(c) == "" {
			name = append(name,c)
		}

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
