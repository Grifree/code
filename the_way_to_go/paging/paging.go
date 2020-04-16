package paging

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Gen struct {
	Page int
	Total int
	PageSize int
	PageCount int
	Url string
}

var errPageCanNotLessZero = errors.New("page can not less zero")
var errTotalCanNotLessZero = errors.New("total can not less zero")
var errPageSizeCanNotLessZero = errors.New("pageSize can not less zero")
var errPageCountCanNotLessZero = errors.New("pageCount can not less zero")
func CreateData(gen Gen) (paging Paging) {
	// 判断错误数据
	if gen.Page <= 0 {
		panic(errPageCanNotLessZero)
	}
	if gen.Total < 0 {
		panic(errTotalCanNotLessZero)
	}
	if gen.PageSize < 0 {
		panic(errPageSizeCanNotLessZero)
	}
	if gen.PageCount < 0 {
		panic(errPageCountCanNotLessZero)
	}
	// 基础数值
	paging.Page = gen.Page
	paging.PageCount = gen.PageCount
	paging.Total = gen.Total
	paging.PageSize = gen.PageSize
	if paging.PageSize == 0 {
		paging.PageSize = 10
	}
	// 计算总页数
	paging.PageCount = getPageCount(paging)
	// TODO page做了修正处理
	if paging.Page > paging.PageCount {
		paging.Page = paging.PageCount
	}
	// 判断是否存在分页
	paging.ExistPaging = getExistPaging(paging)
	paging.Url = getUrl(gen.Url)
	// 计算渲染分页所需数据
	if(paging.ExistPaging){
		paging.getPagingRenderData()
	}
	return paging
}

type Paging struct {
	Page int
	Total int
	PageSize int
	PageCount int
	Url string
	ExistPaging bool
	PagingRenderData
}
type PagingRenderData struct {
	PrevJumpBatch int
	PrevSomePage int
	NextSomePage int
	NextJumpBatch int

	IsFirstPage bool
	IsLastPage bool
	ExistPrevMorePage bool // 存在当前页之前的页码
	ExistPrevBatch bool  // 存在当前页前的"..."
	ExistNextMorePage bool
	ExistNextBatch bool

	PrevPage []int
	NextPage []int
	PrevJumpBatchPage int // "..."对应页数 max(page-7,1)
	NextJumpBatchPage int // "..."对应页数 min(page+7,pageCount)
}

// 以Total PageSize为首选计算依据
func getPageCount(paging Paging) (PageCount int) {
	if paging.Total > 0 {
		PageCount = paging.Total / paging.PageSize
		if(paging.Total % paging.PageSize > 0){
			PageCount += 1
		}
	}
	return PageCount
}

func getExistPaging (paging Paging) (ExistPaging bool) {
	if paging.Page > 0 && paging.PageCount > 0 {
		return true
	}
	return false
}

func getUrl(url string) (pageUrl string) {
	pageUrl = ""
	if strings.Contains(url, "?") {
		pageUrl = strings.Join([]string{url, "&page="}, "")
	} else {
		pageUrl = strings.Join([]string{url, "?page="}, "")
	}
	return pageUrl
}

func (self *Paging) getPagingRenderData() {
	self.PagingRenderData.IsFirstPage = self.Page == 1
	self.PagingRenderData.IsLastPage = self.Page == self.PageCount
	// TODO 改成可传参配置
	self.PagingRenderData.PrevJumpBatch = 7
	self.PagingRenderData.PrevSomePage = 3
	self.PagingRenderData.NextSomePage = 3
	self.PagingRenderData.NextJumpBatch = 7

	if 1 < self.Page {
		self.PagingRenderData.ExistPrevMorePage = true
	}
	if self.Page - 1 > self.PagingRenderData.PrevSomePage {
		self.PagingRenderData.ExistPrevBatch = true
	}
	// 计算 PrevPage
	prevPageLength := 0
	if self.PagingRenderData.ExistPrevBatch {
		prevPageLength = self.PagingRenderData.PrevSomePage
	}else if self.PagingRenderData.ExistPrevMorePage {
		prevPageLength = self.Page - 1
	}
	if prevPageLength > 0 {
		for p := self.Page - prevPageLength; p < self.Page; p++ {
			self.PagingRenderData.PrevPage = append(self.PagingRenderData.PrevPage, p)
		}
	}
	// 计算 PrevJumpBatchPage
	if self.PagingRenderData.ExistPrevBatch {
		prevJumpBatchPage := self.Page - self.PagingRenderData.PrevJumpBatch
		if(prevJumpBatchPage < 1){
			prevJumpBatchPage = 1
		}
		self.PagingRenderData.PrevJumpBatchPage = prevJumpBatchPage
	}

	if self.Page < self.PageCount {
		self.PagingRenderData.ExistNextMorePage = true
	}
	if self.Page + self.PagingRenderData.NextSomePage < self.PageCount {
		self.PagingRenderData.ExistNextBatch = true
	}
	// 计算 NextPage
	nextPageLength := 0
	if self.PagingRenderData.ExistNextBatch {
		nextPageLength = self.PagingRenderData.NextSomePage
	}else if self.PagingRenderData.ExistNextMorePage {
		nextPageLength = self.PageCount - self.Page
	}
	if nextPageLength > 0 {
		for i := 1; i <= nextPageLength; i++ {
			p := self.Page + i
			self.PagingRenderData.NextPage = append(self.PagingRenderData.NextPage, p)
		}
	}
	// 计算 NextJumpBatchPage
	if self.PagingRenderData.ExistNextBatch {
		nextJumpBatchPage := self.Page + self.PagingRenderData.PrevJumpBatch
		if(self.PageCount < nextJumpBatchPage){
			nextJumpBatchPage = self.PageCount
		}
		self.PagingRenderData.NextJumpBatchPage = nextJumpBatchPage
	}
}

// TODO 细节可配置
func (self *Paging) RenderHTML() string  {
	var bufHTML bytes.Buffer
	bufHTML.WriteString("<div>")
	if(self.Total > 0){
		bufHTML.WriteString("<span>总共")
		bufHTML.WriteString(strconv.Itoa(self.Total))
		bufHTML.WriteString("个</span>")
	}

	// 上一页
	if self.PagingRenderData.IsFirstPage {
		bufHTML.WriteString("<span>上一页</span>")
	}else{
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(self.Page - 1))
		bufHTML.WriteString("\">上一页</a>")
	}
	// 第一页 和 ...
	if self.PagingRenderData.ExistPrevBatch {
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString("1\">1</a>")

		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(self.PrevJumpBatchPage))
		bufHTML.WriteString("\">...</a>")
	}
	// 前几页
	for _,p := range self.PrevPage {
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(p))
		bufHTML.WriteString("\">")
		bufHTML.WriteString(strconv.Itoa(p))
		bufHTML.WriteString("</a>")
	}
	// 当前页
	if self.Page > 0 {
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(self.Page))
		bufHTML.WriteString("\">")
		bufHTML.WriteString(strconv.Itoa(self.Page))
		bufHTML.WriteString("</a>")
	}
	// 后几页
	for _,p := range self.NextPage {
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(p))
		bufHTML.WriteString("\">")
		bufHTML.WriteString(strconv.Itoa(p))
		bufHTML.WriteString("</a>")
	}
	// ... 和 最后页
	if self.PagingRenderData.ExistNextBatch {
		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(self.NextJumpBatchPage))
		bufHTML.WriteString("\">...</a>")

		bufHTML.WriteString("<a href=\"")
		bufHTML.WriteString(self.Url)
		bufHTML.WriteString(strconv.Itoa(self.PageCount))
		bufHTML.WriteString("\">")
		bufHTML.WriteString(strconv.Itoa(self.PageCount))
		bufHTML.WriteString("</a>")
	}
	// TODO 跳转指定页

	return bufHTML.String()
}