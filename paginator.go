package utils

import (
	"reflect"
	"sort"
)

var _ Tool = new(PaginatorToolImplement)

var PaginatorTool = new(PaginatorToolImplement)

type PaginatorToolImplement struct {
}

func (p PaginatorToolImplement) SortList(totalList []interface{}, condition *SortCondition) []interface{} {
	if condition == nil {
		return totalList
	}
	if len(totalList) == 0 {
		return totalList
	}
	var newList []interface{}
	var weightList = map[int][]interface{}{}
	var weights []int
	for _, v := range totalList {
		keyWeight := condition.KeyWeightGetter(v, condition.Key)
		if _, exist := weightList[keyWeight]; !exist {
			weightList[keyWeight] = []interface{}{}
		}
		weightList[keyWeight] = append(weightList[keyWeight], v)
		weights = append(weights, keyWeight)
	}
	sort.Ints(weights)
	switch condition.SortType {
	case SortAsc:
		for _, weight := range weights {
			if l, exist := weightList[weight]; exist {
				newList = append(newList, l...)
			}
		}
	case SortDesc:
		for i := len(weights) - 1; i >= 0; i-- {
			if l, exist := weightList[weights[i]]; exist {
				newList = append(newList, l...)
			}
		}
	}
	return newList
}

func (p PaginatorToolImplement) FilterList(list []interface{}, filterConditions []FilterCondition) []interface{} {
	var newList []interface{}
	for _, listItem := range list {
		var suitable = true
		for _, condition := range filterConditions {
			currentValue := condition.KeyValueGetter(listItem, condition.Key)
			switch condition.Operator {
			case In:
				// 判定值是否在指定的范围以内,仅仅适用于判定interface{}是否在[]interface{}之中
				currentValueReflect := reflect.TypeOf(currentValue)
				if currentValueReflect.Kind() == reflect.Slice {
					var inSlice bool
					length := reflect.ValueOf(currentValue).Len()
					for i := 0; i < length; i++ {
						item := reflect.ValueOf(currentValue).Index(i).Interface()
						if reflect.DeepEqual(item, condition.Value) {
							inSlice = true
						}
					}
					if !inSlice {
						// 不在slice里面表示不匹配
						suitable = false
					}
				} else {
					// 范围不是slice直接不匹配
					suitable = false
				}
			case Equals:
				if !reflect.DeepEqual(currentValue, condition.Value) {
					suitable = false
					break
				}
			case NotEquals:
				if reflect.DeepEqual(currentValue, condition.Value) {
					suitable = false
					break
				}
			default:
				// 找不到筛选的条件则默认为符合条件
			}
			if !suitable {
				// 只要有一个filter无法通过则直接跳出循环
				break
			}
		}
		if suitable {
			newList = append(newList, listItem)
		}
	}
	return newList
}

func (p PaginatorToolImplement) GetPage(totalList []interface{}, page uint, limit uint) Paginator {
	var paginator = Paginator{
		Page:     page,
		PageSize: limit,
		Total:    uint(len(totalList)),
		List:     []interface{}{},
	}
	startIndex := (page - 1) * limit
	if startIndex >= paginator.Total {
		// out of range
		return paginator
	}
	endIndex := page * limit
	if endIndex > paginator.Total {
		endIndex = paginator.Total
	}
	paginator.List = totalList[startIndex:endIndex]
	return paginator
}

func (p PaginatorToolImplement) GetFilterPage(totalList []interface{}, page uint, limit uint, filterConditions []FilterCondition, sortCondition *SortCondition) Paginator {

	count := uint(len(totalList))

	if len(filterConditions) != 0 {
		totalList = p.FilterList(totalList, filterConditions)
	}

	if sortCondition != nil {
		totalList = p.SortList(totalList, sortCondition)
	}

	paginator := p.GetPage(totalList, page, limit)
	paginator.Total = count

	return paginator
}
