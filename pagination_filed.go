package pagination

import "net/http"
import "github.com/astaxie/beego/orm"

func (this *Paginator) loadMaxValueByFiled() {
	if n, err := this.omer.LoadRelated(this.Data, this.field); err != nil {
		this.MaxValue = 0
	} else {
		this.MaxValue = n
	}
}
func (this *Paginator) loadCurrentValueByField() {
	this.OffsetFrom = (this.CurrentPage - 1) * this.PerValue
	if this.OffsetFrom > this.MaxValue {
		this.OffsetFrom, this.OffsetTo = 0, 0
		this.CurrentValue = 0
	} else {
		if this.MaxValue >= this.OffsetFrom+this.PerValue {
			this.OffsetTo = this.PerValue
		} else {
			this.OffsetTo = this.MaxValue
		}
		this.CurrentValue = this.OffsetTo - this.OffsetFrom
	}
}

//initilize the basic paginator params
func (p *Paginator) generateFieldData() {
	page := p.Request.URL.Query().Get("page")
	p.CurrentPage = toInt64(page)
	p.loadMaxValueByFiled()
	p.loadMaxPage()
	p.loadCurrentValueByField()
	p.loadPageLinkNext()
	p.loadPageLinkPrev()
	p.loadPageLinkLast()
	p.loadPageLinkFirst()
}

//o: need a query seter, such as o.QueryTable("user")
//tp: the interface to get datas
//params[0], count of show href link in paginatior DEFAULT:5
//params[1], count of show data in per page DEFAULT 15
func NewPaginatorByFiled(r *http.Request, tp interface{}, field string, params ...int64) *Paginator {
	p := &Paginator{
		Request:   r,
		PerValue:  DEFAULT_PERPAGE,
		LinkCount: DEFAULT_LINKCOUNT,
		BasePath:  getBasePath(r.RequestURI),
		Data:      tp,
		omer:      orm.NewOrm(),
		field:     field,
	}

	p.loadQueryArgs()

	//生成数据
	if params != nil && len(params) == 2 {
		p.PerValue = params[1]
	}

	p.generateFieldData()

	//生成分页链接
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
