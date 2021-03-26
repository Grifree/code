package recursive

import "log"

type AreaItem struct {
	ID string
	Name string
	Child []AreaItem
}

var Area = []AreaItem{
	{
		ID:"1",
		Name:"上海",
		Child: []AreaItem{
			{
				ID:"1-1",
				Name:"上海市",
				Child: []AreaItem{
					{
						ID:"1-1-1",
						Name:"宝山区",
						Child: []AreaItem{
							{
								ID:"1-1-1-1",
								Name:"吴淞街道",
							},
							{
								ID:"1-1-1-2",
								Name:"海滨街道",
							},
						},
					},
					{
						ID:"1-1-2",
						Name:"虹口区",
					},
				},
			},
		},
	},
	{
		ID:"2",
		Name:"江苏省",
		Child: []AreaItem{
			{
				ID:"2-1",
				Name:"南京市",
				Child: []AreaItem{
					{
						ID:"2-1-1",
						Name:"玄武区",
					},
					{
						ID:"2-1-2",
						Name:"秦淮区",
					},
					{
						ID:"2-1-3",
						Name:"鼓楼区",
					},
				},
			},
			{
				ID:"2-2",
				Name:"无锡市",
				Child: []AreaItem{
					{
						ID:   "2-2-1",
						Name: "锡山区",
					},
					{
						ID:   "2-2-2",
						Name: "惠山区",
					},
				},
			},
			{
				ID:"2-3",
				Name:"常州市",
			},
		},
	},
}

func FindAllChild(id string) (idList []string) {
	idList = append(idList, id)
	mapTree(&idList,Area)
	return
}

func mapTree(idList *[]string, area []AreaItem) {
	for _,item := range area {
		log.Print(item.ID)
		if isInclude(*idList, item.ID) {
			for _,subItem := range item.Child {
				*idList = append(*idList, subItem.ID)
			}
			mapTree(idList, item.Child)
		}
	}
}

func isInclude(idList []string, id string) (has bool) {
	for _,item := range idList {
		if item == id {
			has = true
			break
		}
	}
	return
}