package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var typeMap = map[string]string{"image/jpeg": ".jpeg",
	"image/jpg": ".jpg",
	"image/png": ".png",
	"image/gif": ".gif",
}

type SliderItem struct {
	LinkUrl string `json:"linkUrl"`
	PicUrl  string `json:"picUrl"`
	Id      int    `json:"id, int"`
}

type Radio struct {
	PicUrl  string `json:"picUrl"`
	Ftitle  string `json:"Ftitle"`
	RadioId int    `json:"id, int"`
}

type Song struct {
	SongName string
}

type Data struct {
	Slider    []SliderItem
	RadioList []Radio
	SongList  []Song
}

type JsonObject struct {
	Code int  `json:"-"`
	Data Data `json:"data"`
}

/*
歌单相关数据类型
*/

type SongList struct {
	Code    int      `json:"-"`
	Data    SongData `json:"data"`
	Default int      `json:"-"`
	Message string   `json:"-"`
	Subcode int      `json:"-"`
}

type SongData struct {
	CategoryId int        `json:"-"`
	Ein        int        `json:"-"`
	List       []SongInfo `json:"list"`
	Sin        int        `json:"-"`
	Sorted     int        `json:"-"`
	Sum        int        `json:"-"`
	Uin        int        `json:"-"`
}

type SongInfo struct {
	Commit_time  string  `json:"-"`
	Createtime   string  `json:"-"`
	Creator      Creator `json:"-"`
	Dissid       string  `json:"-"`
	Dissname     string  `json:"-"`
	Imgurl       string  `json:"imgurl"`
	Introduction string  `json:"-"`
	Listennum    int     `json:"-"`
	Score        int     `json:"-"`
	Version      int     `json:"-"`
}

type Creator struct {
	AvatarUrl   string
	Encrypt_uin string
	Followflag  int
	IsVip       int
	Name        string
	Qq          int
	Type        int
}

func main() {
	//构造http请求

	req, err := http.NewRequest("GET", "https://c.y.qq.com/musichall/fcgi-bin/fcg_yqqhomepagerecommend.fcg", nil)
	if err != nil {
		log.Println(err)
		return
	}

	//设置域名
	req.Host = "c.y.qq.com"
	//设置http头部
	req.Header.Set("Referer", "https://y.qq.com/m/index.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "c.y.qq.com")

	//删除原来目录下的文件
	files, err := ioutil.ReadDir("./vue-music/static")
	if err != nil {
		log.Println(err)
		return
	}

	for _, info := range files {
		os.Remove("./vue-music/static/" + info.Name())
	}

	// 构建文件接收响应体
	file, err := os.OpenFile("./vue-music/static/swiper.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("创建swiper.json文件失败")
		return
	}

	defer file.Close()

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		log.Println("网络请求错误")
		return
	}
	defer response.Body.Close()

	//获取json数据进行分析
	jsonBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		log.Println("将Response.Body读入字节数组时失败")
		return
	}

	//生成一个接收对象
	jsonObject := JsonObject{}

	err = json.Unmarshal(jsonBytes, &jsonObject)
	if err != nil {
		log.Println("解析json数据出错")
		return
	}

	var picUrls = make([]string, 0)
	var picSlider = make([]string, 0)
	var radioUrls = make([]string, 0)
	// var picRadio = make([]string, 0)

	for i := 0; i < 5; i++ {
		picUrls = append(picUrls, jsonObject.Data.Slider[i].PicUrl)
	}

	for i := 0; i < 2; i++ {
		radioUrls = append(radioUrls, jsonObject.Data.RadioList[i].PicUrl)
	}

	//抓取轮播图图片
	for i, url := range picUrls {
		err = fetchImage(url, i+1)
		if err != nil {
			log.Println(err)
			continue
		}
		picSlider = append(picSlider, filepath.Base(url))
		time.Sleep(2 * time.Second)
	}

	//将轮播图图片信息写入json文件
	sliderJson, _ := json.Marshal(picSlider)
	ioutil.WriteFile("./vue-music/static/sliderInfo.json", sliderJson, 0644)

	//抓取电台封面图片
	//for i, url := range radioUrls {
	//	err = fetchImage(url, i)
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//	picRadio = append(picRadio, filepath.Base(url))
	//	time.Sleep(2 * time.Second)
	//}

	////将电台图片信息写入json文件中
	//radioJson, err := json.Marshal(picRadio)
	//ioutil.WriteFile("./vue-music/static/radioInfo.json", radioJson, 0644)

	//抓取歌单数据
	err = fetchSongList()
	if err != nil {
		log.Println("抓取歌单列表失败")
	}

	//抓取歌单列表图片

	//将json文件进行备份
	_, err = file.Write(jsonBytes)
	if err != nil {
		log.Println(err)
		log.Println("文件备份失败")
		return
	}

	log.Println("备份成功")
}

func fetchImage(url string, index int) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	//文件后缀
	base := filepath.Base(url)

	var src string
	if index <= 5 {
		src = "./vue-music/static/" + base
	} else {
		src = "./vue-music/static/" + base
	}

	file, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}

//抓取歌单数据
func fetchSongList() error {
	request, err := http.NewRequest("GET", "https://c.y.qq.com/splcloud/fcgi-bin/fcg_get_diss_by_tag.fcg", nil)
	if err != nil {
		return err
	}

	//设置header
	request.Host = "c.y.qq.com"
	request.Header.Set("Host", "c.y.qq.com")
	request.Header.Set("Referer", "https://c.y.qq.com")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")

	//设置请求参数
	params := request.URL.Query()
	params.Set("categoryId", "10000000")
	params.Set("ein", "29")
	params.Set("sin", "0")
	params.Set("sortId", "5")
	params.Set("platform", "yqq")
	params.Set("loginUin", "0")
	params.Set("hostUin", "0")
	params.Set("g_tk", "5381")
	params.Set("needNewCode", "0")
	rnd := fmt.Sprintf("%f", rand.Float64())
	params.Set("rnd", rnd)
	request.URL.RawQuery = params.Encode()

	client := http.Client{}

	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytesSongList, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result, err := Gbk2Utf8(bytesSongList)
	if err != nil {
		return err
	}

	jsonBytes := ExtractJsonContent(result)

	//获取各封面的URL
	urls := parseSongJson(jsonBytes)

	//抓取封面图片,并将封面图片信息写入json文件中
	covers := make([]string, 0)
	for i, url := range urls {
		cover, err := fetchCover(url, i)
		if err != nil {
			log.Println(err)
		}
		covers = append(covers, cover)
	}

	jsonCovers, err := json.Marshal(&covers)
	if err != nil {
		return err
	}

	//保存songList的封面信息
	err = ioutil.WriteFile("./vue-music/static/songListCover.json", jsonCovers, 0644)
	if err != nil {
		return err
	}

	//保存songList信息
	err = ioutil.WriteFile("./vue-music/static/songList.json", jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

//gbk转码为UTF-8
func Gbk2Utf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//正则抽取json内容
func ExtractJsonContent(s []byte) []byte {
	r := regexp.MustCompile(`^MusicJsonCallback\((.*?)\)$`)
	params := r.FindSubmatch(s)
	return params[1]
}

//解析songList的json文件,返回封面URL的数组
func parseSongJson(jsonbytes []byte) []string {
	var jsonObj = SongList{}
	err := json.Unmarshal(jsonbytes, &jsonObj)
	if err != nil {
		log.Println("解析歌单json数据出错")
	}

	songlist := jsonObj.Data.List

	result := make([]string, 0)

	for _, val := range songlist {
		result = append(result, val.Imgurl)
	}
	return result
}

//抓取封面函数
func fetchCover(url string, index int) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var ext string
	t := res.Header.Get("Content-Type")
	for k, v := range typeMap {
		if t == k {
			ext = v
			break
		}
	}

	if ext == "" {
		return "", errors.New("不存在指定类型的图片")
	}

	n := strconv.Itoa(index)

	filename := "cover_" + n + ext

	file, err := os.OpenFile("./vue-music/static/"+filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}
