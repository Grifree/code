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
	intervalFirstToCurrentPage := 0 // 第一页和当前页之间有几个页
	if !render.IsFirstPage {
		intervalFirstToCurrentPage = naturalNumberInterval(1, gen.Page)
	}
	closestPrevPageLength := 0
	if intervalFirstToCurrentPage > 0 {
		closestPrevPageLength = minInt(intervalFirstToCurrentPage, gen.ClosestPageLength)
		render.ClosestPage.Prev.Exist = closestPrevPageLength > 0
	}
	// ClosestPage.Prev.PageList
	if render.ClosestPage.Prev.Exist {
		for i:=1; i<=closestPrevPageLength; i++ {
			curSlice := []int{gen.Page - i}
			render.ClosestPage.Prev.PageList = append(curSlice, render.ClosestPage.Prev.PageList...)
		}
	}

	// ClosestPage.Next.Exist
	intervalCurrentToLastPage := 0 // 当前页和最后一页之间有几个页
	if !render.IsLastPage {
		intervalCurrentToLastPage = naturalNumberInterval(gen.Page, render.LastPage)
	}
	closestNextPageLength := 0
	if intervalCurrentToLastPage > 0 {
		closestNextPageLength = minInt(intervalCurrentToLastPage, gen.ClosestPageLength)
		render.ClosestPage.Next.Exist = closestNextPageLength > 0
	}
	// ClosestPage.Next.PageList
	if render.ClosestPage.Next.Exist {
		for i:=1; i<=closestNextPageLength; i++ {
			render.ClosestPage.Next.PageList = append(render.ClosestPage.Next.PageList, gen.Page + i)
		}
	}

	// JumpPage.Prev
	render.JumpPage.Prev.Exist = intervalFirstToCurrentPage > gen.ClosestPageLength
	if render.JumpPage.Prev.Exist {
		render.JumpPage.Prev.Interval = maxInt(1, gen.Page - gen.JumpPageInterval)
	}

	// JumpPage.Next
	render.JumpPage.Next.Exist = intervalCurrentToLastPage > gen.ClosestPageLength
	if render.JumpPage.Next.Exist {
		render.JumpPage.Next.Interval = minInt(gen.Page + gen.JumpPageInterval, render.LastPage)
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
func minInt(a int, b int) (maxInt int) {
	if(a < b){
		return a
	}
	return b
}
// TODO 测试
func maxInt(a int, b int) (maxInt int) {
	if(a < b){
		return b
	}
	return a
}
// TODO 测试
func naturalNumberInterval(a int, b int) (interval int) {
	var min, max int
	if a > b {
		min = b
		max = a
	} else if a < b {
		min = a
		max = b
	} else {
		// a == b
		return 0
	}
	// 1 2 3 4 5 6 7 8 9
	// min 1 max 5
	// 5-1-1
	// min 5 max 9
	// 9-5-1
	return max - min - 1
}