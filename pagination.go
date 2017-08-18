package pagination

import "net/http"
import "strings"
import "github.com/astaxie/beego/orm"
import "strconv"
import "fmt"

const (
	DEFAULT_PERPAGE   = 15
	DEFAULT_LINKCOUNT = 5
)

type Link struct {
	Id        int64
	IsCurrent bool
	Href      string
}

type Paginator struct {
	Request *http.Request

	PageLinkPrev  string
	PageLinkNext  string
	PageLinkLast  string
	PageLinkFirst string

	PerValue     int64
	MaxValue     int64
	CurrentValue int64

	CurrentPage int64
	MaxPage     int64

	BasePath  string
	Data      interface{}
	Links     []Link
	LinkCount int64

	o orm.QuerySeter
}

//generate url include query string
func LoadLink(path, key string, value interface{}) string {
	var url string
	if strings.Index(path, "?") == -1 {
		url = fmt.Sprintf("%s?%s=", path, key)
	} else {
		url = fmt.Sprintf("%s&%s=", path, key)
	}

	switch v := value.(type) {
	case int:
		url += strconv.Itoa(v)
	case int64:
		url += strconv.Itoa(int(v))
	case string:
		url += v
	}

	return url
}

//get the base path, such as /a/b
func getBasePath(url string) string {
	if n := strings.Index(url, "?"); n != -1 {
		return url[:n]
	}
	return url
}

//convert string to int64
func toInt64(str string) int64 {
	if n, err := strconv.Atoi(str); err != nil {
		return 1
	} else if n < 0 {
		return 1
	} else {
		return int64(n)
	}

}

func (p *Paginator) loadMaxValue() {
	if n, err := p.o.Count(); err != nil {
		p.MaxValue = 0
	} else {
		p.MaxValue = n
	}
}

func (p *Paginator) loadMaxPage() {
	if p.MaxValue%p.PerValue == 0 {
		p.MaxPage = p.MaxValue / p.PerValue
		return
	}
	p.MaxPage = p.MaxValue/p.PerValue + 1
}

func (p *Paginator) loadData() {
	p.CurrentValue, _ = p.o.Offset((p.CurrentPage - 1) * p.PerValue).Limit(p.PerValue).All(p.Data)
}

func (p *Paginator) loadLinks() {
	var links []Link
	if p.CurrentPage+p.LinkCount > p.MaxPage {
		links = p.generateLinks(p.MaxPage-p.LinkCount+1, p.MaxPage)
	} else if p.CurrentPage%p.LinkCount == 1 {
		links = p.generateLinks(p.CurrentPage, p.CurrentPage+p.LinkCount-1)
	} else {
		index := (p.CurrentPage - 1) / p.LinkCount
		start := index*p.LinkCount + 1
		links = p.generateLinks(start, start+p.LinkCount-1)
	}

	p.Links = links
}

func (p *Paginator) loadPageLinkNext() {
	if p.CurrentPage+1 > p.MaxPage {
		p.PageLinkNext = "#"
		return
	}
	p.PageLinkNext = LoadLink(p.BasePath, "page", p.CurrentPage+1)
}

func (p *Paginator) loadPageLinkPrev() {
	if p.CurrentPage-1 < 1 {
		p.PageLinkPrev = "#"
		return
	}
	p.PageLinkPrev = LoadLink(p.BasePath, "page", p.CurrentPage-1)
}

func (p *Paginator) loadPageLinkFirst() {
	p.PageLinkFirst = LoadLink(p.BasePath, "page", 1)
}

func (p *Paginator) loadPageLinkLast() {
	p.PageLinkLast = LoadLink(p.BasePath, "page", p.MaxPage)
}

//initilize the basic paginator params
func (p *Paginator) generateData() {
	page := p.Request.URL.Query().Get("page")
	p.CurrentPage = toInt64(page)
	p.loadMaxValue()
	p.loadMaxPage()
	p.loadData()
	p.loadPageLinkNext()
	p.loadPageLinkPrev()
	p.loadPageLinkLast()
	p.loadPageLinkFirst()
}

//initilize href
func (p *Paginator) generateLinks(from, to int64) []Link {
	var links []Link
	for i := from; i <= to; i++ {
		links = append(links, Link{
			Id:        i,
			IsCurrent: i == p.CurrentPage,
			Href:      LoadLink(p.BasePath, "page", i),
		})
	}

	return links
}

func NewPaginator(r *http.Request, o orm.QuerySeter, tp interface{}, params ...int64) *Paginator {
	p := &Paginator{
		Request:   r,
		PerValue:  DEFAULT_PERPAGE,
		LinkCount: DEFAULT_LINKCOUNT,
		BasePath:  getBasePath(r.RequestURI),
		Data:      tp,
		o:         o,
	}

	if params != nil && len(params) == 2 {
		p.PerValue = params[1]
	}
	p.generateData()

	if params != nil {
		if params[0] > p.MaxPage {
			p.LinkCount = p.MaxPage
		} else {
			p.LinkCount = params[0]
		}
	}
	p.loadLinks()

	return p
}
