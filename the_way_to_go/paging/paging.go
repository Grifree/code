package paging

import (
	"errors"
	"log"
)

type Gen struct {
	Page int
	Total int
	PerPage int
	JumpPageInterval int
	ClosestPageLength int
}
type Render struct {
	ExistPaging bool

	LastPage int
	IsFirstPage bool
	IsLastPage bool

	ClosestPage struct {
		Prev struct{
			Exist bool
			PageList []int // TODO 返回时 零值设为空数组,而非nil
		}
		Next struct{
			Exist bool
			PageList []int
		}
	}

	JumpPage struct{
		Prev struct{
			Exist bool
			Interval int
		}
		Next struct{
			Exist bool
			Interval int
		}
	}
}

func CreateData(gen Gen) (Render, Gen) {
	genCheckAndFix(&gen)

	render := Render{}
	render.ExistPaging = gen.Total > 0
	if !render.ExistPaging {
		return render, gen
	}
	// LastPage
	render.LastPage = gen.Total / gen.PerPage
	if gen.Total % gen.PerPage > 0 {
		render.LastPage += 1
	}
	// 修正过大Page
	if gen.Page > render.LastPage {
		gen.Page = render.LastPage
	}

	render.IsFirstPage = gen.Page == 1
	render.IsLastPage = gen.Page == render.LastPage

	// ClosestPage.Prev.Exist
	pageCountBetweenFirstToCurrentPage := 0 // 第一页和当前页之间有几个页
	if !render.IsFirstPage {
		pageCountFromFirstToCurrentPage := gen.Page - 1
		pageCountBetweenFirstToCurrentPage = pageCountFromFirstToCurrentPage - 1
	}
	closestPagePrevPageLength := 0
	if pageCountBetweenFirstToCurrentPage > 0 {
		closestPagePrevPageLength = getIntMin(pageCountBetweenFirstToCurrentPage, gen.ClosestPageLength)
		render.ClosestPage.Prev.Exist = closestPagePrevPageLength > 0
	}
	// ClosestPage.Prev.PageList
	if render.ClosestPage.Prev.Exist {
		for i:=1; i<=closestPagePrevPageLength; i++ {
			curSlice := []int{gen.Page - i}
			render.ClosestPage.Prev.PageList = append(curSlice, render.ClosestPage.Prev.PageList...)
		}
	}

	// ClosestPage.Next.Exist
	pageCountBetweenCurrentToLastPage := 0 // 当前页和最后一页之间有几个页
	if !render.IsLastPage {
		pageCountFromCurrentToLastPage := render.LastPage - gen.Page
		pageCountBetweenCurrentToLastPage = pageCountFromCurrentToLastPage - 1
	}
	closestPageNextPageLength := 0
	if pageCountBetweenCurrentToLastPage > 0 {
		closestPageNextPageLength = getIntMin(pageCountBetweenCurrentToLastPage, gen.ClosestPageLength)
		render.ClosestPage.Next.Exist = closestPageNextPageLength > 0
	}
	// ClosestPage.Next.PageList
	if render.ClosestPage.Next.Exist {
		for i:=1; i<=closestPageNextPageLength; i++ {
			render.ClosestPage.Next.PageList = append(render.ClosestPage.Next.PageList, gen.Page + i)
		}
	}

	// JumpPage.Prev
	render.JumpPage.Prev.Exist = pageCountBetweenFirstToCurrentPage > gen.ClosestPageLength
	if render.JumpPage.Prev.Exist {
		render.JumpPage.Prev.Interval = getIntMax(1, gen.Page - gen.JumpPageInterval)
	}

	// JumpPage.Next
	render.JumpPage.Next.Exist = pageCountBetweenCurrentToLastPage > gen.ClosestPageLength
	if render.JumpPage.Next.Exist {
		render.JumpPage.Next.Interval = getIntMin(gen.Page + gen.JumpPageInterval, render.LastPage)
	}

	return render, gen
}

const errMsgPrefix = "paging: paging.CreateData(gen) "
var errPageCannotLessZero = errors.New(errMsgPrefix+"gen.Page cannot less zero")
var errPageCannotBeZeroAndFix = errors.New(errMsgPrefix+"gen.Page can not be 0, it's set to 1, but you need check your code")
var errTotalCannotLessZero = errors.New(errMsgPrefix+"gen.Total cannot less zero")
var errPerPageCannotLessZero = errors.New(errMsgPrefix+"gen.PerPage cannot less zero")
var errPerPageCannotBeZero = errors.New(errMsgPrefix+"gen.PerPage cannot be 0")
var errJumpPageIntervalCannotLessZero = errors.New(errMsgPrefix+"gen.JumpPageInterval cannot less zero")
var errJumpPageIntervalCannotBeZero = errors.New(errMsgPrefix+"gen.JumpPageInterval cannot be 0")
var errClosestPageLengthCannotLessZero = errors.New(errMsgPrefix+"gen.ClosestPageLength cannot less zero")
func genCheckAndFix (genPtr *Gen) (pass bool){
	if genPtr.Page < 0 {
		panic(errPageCannotLessZero)
	}
	if genPtr.Page == 0 {
		genPtr.Page = 1
		log.Print(errPageCannotBeZeroAndFix)
	}
	if genPtr.Total < 0 {
		panic(errTotalCannotLessZero)
	}
	if genPtr.PerPage < 0 {
		panic(errPerPageCannotLessZero)
	}
	if genPtr.PerPage == 0 {
		panic(errPerPageCannotBeZero)
	}
	if genPtr.JumpPageInterval < 0 {
		panic(errJumpPageIntervalCannotLessZero)
	}
	if genPtr.JumpPageInterval == 0 {
		panic(errJumpPageIntervalCannotBeZero)
	}
	if genPtr.ClosestPageLength < 0 {
		panic(errClosestPageLengthCannotLessZero)
	}
	// TODO pass 怎么测试
	return true
}

// TODO 测试
func getIntMin(a int, b int) (maxInt int) {
	if(a < b){
		return a
	}
	return b
}
// TODO 测试
func getIntMax(a int, b int) (maxInt int) {
	if(a < b){
		return b
	}
	return a
}