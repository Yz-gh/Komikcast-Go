package komikcast

import (
    "fmt"
    "time"
    "strings"
    "encoding/json"
    
    "github.com/valyala/fasthttp"
)

type ComicInfo struct {
	Linkid       string   `json:"linkid"`
	Title        string   `json:"title"`
	TitleOther   string   `json:"title_other"`
	Author       string   `json:"author"`
	Image        string   `json:"image"`
	Image2       string   `json:"image2"`
	Rating       string   `json:"rating"`
	Sinopsis     string   `json:"sinopsis"`
	Type         string   `json:"type"`
	Status       string   `json:"status"`
	Released     string   `json:"released"`
	TotalChapter string   `json:"total_chapter"`
	UpdatedOn    string   `json:"updated_on"`
	Genres       []string `json:"genres"`
	ListChapter  []*Chapter `json:"list_chapter"`
}

type ChapterDetail struct {
	Title       string `json:"title"`
	Ch          string `json:"ch"`
	ComicTitle  string `json:"comic_title"`
	PrevCh      string `json:"prev_ch"`
	NextCh      string `json:"next_ch"`
	PrevLinkID  string `json:"prev_link_id"`
	NextLinkID  string `json:"next_link_id"`
	ListChapter []*Chapter `json:"list_chapter"`
	Images []string `json:"images"`
}

type ChapterByPage struct {
	CurrentPage int  `json:"currentPage"`
	PerPage     int  `json:"perPage"`
	Total       int  `json:"total"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
	TotalPages  int  `json:"totalPages"`
	Data        []*Chapter `json:"data"`
}

type Chapter struct {
	Title  string `json:"title"`
	Image  string `json:"image"`
	Image2 string `json:"image2"`
	IsHot  string `json:"isHot"`
	Link   string `json:"link"`
	LinkID string `json:"linkId"`
	Ch     string `json:"ch"`
	Chapter    string `json:"chapter"`
	TimeRelease string `json:"time_release"`
	ChID   string `json:"ch_id"`
	ChTime string `json:"ch_time"`
	IsCompleted string `json:"isCompleted"`
	Type string `json:"type"`
	Rating string `json:"rating"`
}

var host = "https://apk.nijisan.my.id"
type RequestDetail struct {
	Method   string   `json:"method"`
	Endpoint string   `json:"endpoint"`
	Path     []string `json:"path"`
	Query    []string `json:"query"`
}
    
var path = map[string]*RequestDetail{
    "read_and_info": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/baca/%s",
        Path: []string{"id"}},
    "filter": &RequestDetail{
        Method: "POST",
        Endpoint: "/premium/komik/filter"},
    "chapterList": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/info/%s/ch",
        Path: []string{"id"},
        Query: []string{"page", "limit"}},
    "genre": &RequestDetail{
        Method: "GET",
        Endpoint: "/filter/komik/genre"},
    /*"randomGenre": &RequestDetail{
        Method: "GET",
        Endpoint: "/filter/komik/genre/random"},*/
    "home": &RequestDetail{
        Method: "GET",
        Endpoint: "/premium/home"},
    "latest": &RequestDetail{
        Method: "GET",
        Endpoint: "/premium/home/latest/%s/%s",
        Path: []string{"newpage", "page"}},
    "recommended": &RequestDetail{
        Method: "GET",
        Endpoint: "/filter/komik/rekomendasi/%s/%s",
        Path: []string{"type", "limit"}},
    "search": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/search/%s",
        Path: []string{"keyword"}},
    "searchByPage": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/search/%s/%s/%s",
        Path: []string{"keyword", "newpage", "page"}},
}

func addParams(url string, k, v []string) string{
    if len(k) < len(v){ return url }
    url += "?"
    p := ""
    for i := range v{
        if i != 0 { p += "&" }
        p += k[i] + "=" + v[i]
    }
    return url+p
}

func request(url, method string, jsonBody []byte) string{
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    req.Header.SetMethod(method)
    req.SetRequestURI(url)
    if method == "POST"{
      req.SetBody(jsonBody)
      req.Header.SetContentType("application/json")
    }
    resp := fasthttp.AcquireResponse()
    if err := fasthttp.DoTimeout(req, resp, 10*time.Second); err != nil {
      return err.Error()
    }
    defer fasthttp.ReleaseResponse(resp)
    return string(resp.Body())
}

func readOrInfo(id string) (c interface{}, err error){
    rd := path["read_and_info"]
    url := fmt.Sprintf(host+rd.Endpoint, id)
    respBody := request(url, rd.Method, nil)
    if strings.Contains(id, "chapter"){ c = &ChapterDetail{} } else { c = &ComicInfo{} }
    err = json.Unmarshal([]byte(respBody), c)
    return
}

func ReadComic(id string) (*ChapterDetail, error){
    r, err := readOrInfo(id)
    return r.(*ChapterDetail), err
}

func GetComic(id string) (*ComicInfo, error){
    r, err := readOrInfo(id)
    return r.(*ComicInfo), err
}

func GetChapterByPage(id, page, limit string) (c *ChapterByPage, err error){
    rd := path["chapterList"]
    url := fmt.Sprintf(host+rd.Endpoint, id)
    url = addParams(url, rd.Query, []string{page, limit})
    respBody := request(url, rd.Method, nil)
    c = &ChapterByPage{}
    err = json.Unmarshal([]byte(respBody), c)
    return
}

func GenreList() (gl []string, err error){
    rd := path["genre"]
    url := host+rd.Endpoint
    respBody := request(url, rd.Method, nil)
    r := make(map[string][]string)
    err = json.Unmarshal([]byte(respBody), &r)
    gl = r["genre"]
    return
}

func Home() (lc map[string][]*Chapter, err error){
    rd := path["home"]
    url := host+rd.Endpoint
    respBody := request(url, rd.Method, nil)
    err = json.Unmarshal([]byte(respBody), &lc)
    return
}

func GetLatestUpdate(page string) (lc []*Chapter, err error){
    rd := path["latest"]
    url := fmt.Sprintf(host+rd.Endpoint, "1", page)
    respBody := request(url, rd.Method, nil)
    r := make(map[string]interface{})
    err = json.Unmarshal([]byte(respBody), &r)
    ld := r["data"].([]interface{})
    for _, data := range ld{ 
        ta := &Chapter{}
        dataByte, _ := json.Marshal(data.(map[string]interface{}))
        json.Unmarshal(dataByte, ta)
        lc = append(lc, ta)
    }
    return
}

func GetRecommendedComic(tipe, limit string) (lc []*Chapter, err error){
    rd := path["recommended"]
    url := fmt.Sprintf(host+rd.Endpoint, tipe, limit)
    respBody := request(url, rd.Method, nil)
    err = json.Unmarshal([]byte(respBody), &lc)
    return
}

func SearchComic(keyword string) (lc []*Chapter, err error){
    rd := path["search"]
    url := fmt.Sprintf(host+rd.Endpoint, keyword)
    respBody := request(url, rd.Method, nil)
    /*r := make(map[string][]*Chapter)
    err = json.Unmarshal([]byte(respBody), &r)
    lc = r["page"]*/
    r := make(map[string]interface{})
    err = json.Unmarshal([]byte(respBody), &r)
    lp := r["page"].([]interface{})
    for _, page := range lp{ 
        ta := &Chapter{}
        pageByte, _ := json.Marshal(page.(map[string]interface{}))
        json.Unmarshal(pageByte, ta)
        lc = append(lc, ta)
    }
    return
}
