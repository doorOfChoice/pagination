# Paginator

It's can generate any href you want, and it depends on beego/orm

## Install

go get github.com/doorOfChoice/pagination

## Use
```go
    //o: need a query seter, such as o.QueryTable("user")
    //tp: the interface to get datas
    //params[0], count of show href link in paginatior DEFAULT:5
    //params[1], count of show data in per page DEFAULT 15
    func NewPaginator(r *http.Request, o orm.QuerySeter, tp interface{}, params ...int64) *Paginator 
```

userController.go
```go
    // If I has a model is User
    type User struct{
        Id int `orm:"auto;pk"`
        Name string 
    }

    func (c *UserController) GetAllUsers() {
        var users []*User

        o := orm.NewOrm()
        //generate a paginator
        paginator := pagination.NewPaginator(
            c.Ctx.Request,
            o.Query("user").Filter("id" > 5),
            &users,
            15,
            15,
        )

        c.Data["P"] = paginator
        c.Data["Users"] = users
        c.TplName = "user.tpl"
    }
```

## Param

There are some params, you can use it in template
```go
   type Paginator struct {

	PageLinkPrev  string 
	PageLinkNext  string
	PageLinkLast  string
	PageLinkFirst string

	PerValue     int64
	MaxValue     int64  //sum of data in the table
	CurrentValue int64  //how many data in current page

	CurrentPage int64
	MaxPage     int64

	BasePath  string
	Links     []Link    //generated links
}

type Link struct {
	Id        int64 //which page
	IsCurrent bool  //is current page
	Href      string//the href
}
```

## License

MIT