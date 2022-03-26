package komikcast

import (
    "fmt"
    "time"
    "strings"
    "encoding/json"
    
    "github.com/valyala/fasthttp"
)

type ComicInfo struct {
	Linkid       string   `json:"linkid,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleOther   string   `json:"title_other,omitempty"`
	Author       string   `json:"author,omitempty"`
	Image        string   `json:"image,omitempty"`
	Image2       string   `json:"image2,omitempty"`
	Rating       string   `json:"rating,omitempty"`
	Sinopsis     string   `json:"sinopsis,omitempty"`
	Type         string   `json:"type,omitempty"`
	Status       string   `json:"status,omitempty"`
	Released     string   `json:"released,omitempty"`
	TotalChapter string   `json:"total_chapter,omitempty"`
	UpdatedOn    string   `json:"updated_on,omitempty"`
	Genres       []string `json:"genres,omitempty"`
	ListChapter  []*Chapter `json:"list_chapter,omitempty"`
}

type ChapterDetail struct {
	Title       string `json:"title,omitempty"`
	Ch          string `json:"ch,omitempty"`
	ComicTitle  string `json:"comic_title,omitempty"`
	PrevCh      string `json:"prev_ch,omitempty"`
	NextCh      string `json:"next_ch,omitempty"`
	PrevLinkID  string `json:"prev_link_id,omitempty"`
	NextLinkID  string `json:"next_link_id,omitempty"`
	ListChapter []*Chapter `json:"list_chapter,omitempty"`
	Images []string `json:"images,omitempty"`
}

type ChapterByPage struct {
	CurrentPage int  `json:"currentPage,omitempty"`
	PerPage     int  `json:"perPage,omitempty"`
	Total       int  `json:"total,omitempty"`
	HasNextPage bool `json:"hasNextPage,omitempty"`
	HasPrevPage bool `json:"hasPrevPage,omitempty"`
	TotalPages  int  `json:"totalPages,omitempty"`
	Data        []*Chapter `json:"data,omitempty"`
}

type Chapter struct {
	Title  string `json:"title,omitempty"`
	Image  string `json:"image,omitempty"`
	Image2 string `json:"image2,omitempty"`
	IsHot  string `json:"isHot,omitempty"`
	Link   string `json:"link,omitempty"`
	LinkID string `json:"linkId,omitempty"`
	Ch     string `json:"ch,omitempty"`
	Chapter    string `json:"chapter,omitempty"`
	TimeRelease string `json:"time_release,omitempty"`
	ChID   string `json:"ch_id,omitempty"`
	ChTime string `json:"ch_time,omitempty"`
	IsCompleted string `json:"isCompleted,omitempty"`
	Type string `json:"type,omitempty"`
	Rating string `json:"rating,omitempty"`
}

type ResultFilter struct {
	Genre        []string `json:"genre,omitempty"`
	Status       string   `json:"status,omitempty"`
	Order        string   `json:"order,omitempty"`
	FilterResult []*Chapter `json:"filter_result,omitempty"`
}

type SearchResult struct {
	Keyword string `json:"keyword"`
	Page    []*Chapter `json:"page"`
}

var host = "https://apk.nijisan.my.id"
type RequestDetail struct {
	Method   string
	Endpoint string
	Path     []string
	Query    []string
}
    
var path = map[string]*RequestDetail{
    "read_and_info": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/baca/%s",
        Path: []string{"id"}},
    "filter": &RequestDetail{
        Method: "POST",
        Endpoint: "/komik/filter"},
    "chapterList": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/info/%s/ch",
        Path: []string{"id"},
        Query: []string{"page", "limit"}},
    "genre": &RequestDetail{
        Method: "GET",
        Endpoint: "/komik/genre"},
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
        Endpoint: "/komik/rekomendasi/%s/%s",
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

func ReadOrInfo(id string) (c interface{}, err error){
    rd := path["read_and_info"]
    url := fmt.Sprintf(host+rd.Endpoint, id)
    respBody := request(url, rd.Method, nil)
    if strings.Contains(respBody, "prev"){ c = &ChapterDetail{} } else { c = &ComicInfo{} }
    err = json.Unmarshal([]byte(respBody), c)
    return
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

func SearchComicV2(keyword, page string) (r *SearchResult, err error){
    rd := path["searchByPage"]
    url := fmt.Sprintf(host+rd.Endpoint, keyword, "1", page)
    respBody := request(url, rd.Method, nil)
    r = &SearchResult{}
    err = json.Unmarshal([]byte(respBody), &r)
    return
}

/* Filter Comic Args Exmp
  sort = title, titlereverse, latest, popular
  status = Ongoing, Completed, ""(blank) if want all
  genre = see GenreList()
*/
func FilterComic(page, status, order string, genres []string) (rf *ResultFilter, err error){
    rd := path["filter"]
    url := fmt.Sprintf(host+rd.Endpoint)
    g := ""
    if len(genres) != 0{
        gs, _ := json.Marshal(genres)
        g = `"genre": `+string(gs)+`,`
    }
    jsonToSend := fmt.Sprintf(`{
      %s
      "newpage": 1,
      "page": %s,
      "status": "%s",
      "order": "%s"
    }`, g, page, status, order)
    respBody := request(url, rd.Method, []byte(jsonToSend))
    rf = &ResultFilter{}
    err = json.Unmarshal([]byte(respBody), rf)
    return
}
