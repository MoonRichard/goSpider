package Rpider_v1

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
	urlf := `http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2016/11/01/110118.html`
	l := `<a href='18/110118001.html'>110118001000</a></td><td><a href='18/110118001.html'>鼓楼街道办事处</a></td></tr><tr class='towntr'><td><a href='18/110118002.html'>110118002000</a></td><td><a href='18/110118002.html'>果园街道办事处</a></td></tr><tr class='towntr'><td><a href='18/110118003.html'>110118003000</a></td><td><a href='18/110118003.html'>檀营地区办事处</a></td></tr><tr class='towntr'><td><a href='18/110118100.html'>110118100000</a></td><td><a href='18/110118100.html'>密云镇</a></td></tr><tr class='towntr'><td><a href='18/110118101.html'>110118101000</a></td><td><a href='18/110118101.html'>溪翁庄镇</a></td></tr><tr class='towntr'><td><a href='18/110118102.html'>110118102000</a></td><td><a href='18/110118102.html'>西田各庄镇</a></td></tr><tr class='towntr'><td><a href='18/110118103.html'>110118103000</a></td><td><a href='18/110118103.html'>十里堡镇</a></td></tr><tr class='towntr'><td><a href='18/110118104.html'>110118104000</a></td><td><a href='18/110118104.html'>河南寨镇</a></td></tr><tr class='towntr'><td><a href='18/110118105.html'>110118105000</a></td><td><a href='18/110118105.html'>巨各庄镇</a></td></tr><tr class='towntr'><td><a href='18/110118106.html'>110118106000</a></td><td><a href='18/110118106.html'>穆家峪镇</a></td></tr><tr class='towntr'><td><a href='18/110118107.html'>110118107000</a></td><td><a href='18/110118107.html'>太师屯镇</a></td></tr><tr class='towntr'><td><a href='18/110118108.html'>110118108000</a></td><td><a href='18/110118108.html'>高岭镇</a></td></tr><tr class='towntr'><td><a href='18/110118109.html'>110118109000</a></td><td><a href='18/110118109.html'>不老屯镇</a></td></tr><tr class='towntr'><td><a href='18/110118110.html'>110118110000</a></td><td><a href='18/110118110.html'>冯家峪镇</a></td></tr><tr class='towntr'><td><a href='18/110118111.html'>110118111000</a></td><td><a href='18/110118111.html'>古北口镇</a></td></tr><tr class='towntr'><td><a href='18/110118112.html'>110118112000</a></td><td><a href='18/110118112.html'>大城子镇</a></td></tr><tr class='towntr'><td><a href='18/110118113.html'>110118113000</a></td><td><a href='18/110118113.html'>东邵渠镇</a></td></tr><tr class='towntr'><td><a href='18/110118114.html'>110118114000</a></td><td><a href='18/110118114.html'>北庄镇</a></td></tr><tr class='towntr'><td><a href='18/110118115.html'>110118115000</a></td><td><a href='18/110118115.html'>新城子镇</a></td></tr><tr class='towntr'><td><a href='18/110118116.html'>110118116000</a></td><td><a href='18/110118116.html'>石城镇</a></td></tr><tr class='towntr'><td><a href='18/110118400.html'>110118400000</a></td><td><a href='18/110118400.html'>北京密云经济开发区</a></td></tr>
</table>`
	list := getUrlList(l)
	for _,l := range list {
		parsePage(urlf,l,"++")
	}
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
			parsePage(baseURL,urlList[i],"  ")
		}
	}else {
		fmt.Println("解析首页错误")
	}
}

func parsePage(urlFather string,sub string,outPrefix string)(idList []string,nameList []string,urlList []string) {
	url := mixUrl(urlFather,sub)
	fmt.Println(url)
	resp,_ := http.Get(url)
	body,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//转码了解一下烦死了
	utf8 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	table := getList(utf8)
	idList , nameList = getIdAndName(table)
	urlList = getUrlList(table)
	i := 0
	for ; i<len(idList);i++  {
		count++
		fmt.Print(outPrefix+idList[i]+" "+nameList[i])
		fmt.Println(""+strconv.Itoa(count))
		parsePage(url,urlList[i],outPrefix+"  ")
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
