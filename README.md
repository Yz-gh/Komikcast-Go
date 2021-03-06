# Komikcast-Go Example
```go
package main

import (
    "fmt"

    "github.com/Yz-gh/Komikcast-Go"
)

func main(){
    fmt.Println("=====[ Search ]======")
    s, _ := komikcast.SearchComic("omniscient")
    for _, cd := range s{ fmt.Printf("%+v\n", cd) }
    cid := s[0].LinkID
    fmt.Println()

    fmt.Println("=====[ Get Comic ]=====")
    fmt.Println("Comic id to get:", cid)
    c, _ := komikcast.GetComic(cid)
    ca := c.(*komikcast.ComicInfo)
    fmt.Printf("%+v\nChapter list:\n", ca)
    for _, d := range ca.ListChapter{ fmt.Printf("%+v\n", d) }
    fmt.Println()

    fmt.Println("=====[ Read Comic ]=====")
    tr := ca.ListChapter[0].LinkID
    fmt.Println("Chapter id to read:", tr)
    rc, _ := komikcast.ReadComic(tr)
    rca := rc.(*komikcast.ChapterDetail)
    for _, img := range rca.Images{ fmt.Println(img) }
    fmt.Println()

    fmt.Println("=====[ Filter Comic ]=====")
    fc, _ := komikcast.FilterComic("1", "", "latest", []string{"Super Power", "Reincarnation"})
    for _, f := range fc.FilterResult{ fmt.Printf("%+v\n", f) }
}
