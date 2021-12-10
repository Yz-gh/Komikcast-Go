package komikcast

import (
    "fmt"
    "time"
    "strings"
    "encoding/json"
    
    "github.com/valyala/fasthttp"
)

type Chapter struct {
	Ch          string `json:"ch"`
	TimeRelease string `json:"time_release"`
	LinkID      string `json:"linkId"`
}

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

var host = ""
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
