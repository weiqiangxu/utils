package utils

import (
	"fmt"
	"strconv"
	"testing"
)

type School struct {
	Name string
	Info struct {
		Age     int
		Leader  string
		Teacher []string
	}
}

func TestPaginatorToolImplement_SortList(t *testing.T) {

	// 自定义权重排位并且权重以整数呈现
	var keyWeightGetter = func(listItem interface{}, key string) int {
		if school, ok := listItem.(School); ok {
			switch key {
			case "age":
				return school.Info.Age
			case "leader":
				// 如果是字符串，那么将字符串ASCII码作为权重
				if school.Info.Leader == "" {
					return 0
				}
				i, err := strconv.Atoi(fmt.Sprintf("%d", school.Info.Leader[0]))
				if err != nil {
					return 0
				}
				return i
			}

		}
		return 0
	}

	var list []School
	list = append(list, School{
		Name: "广州",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:    18,
			Leader: "A",
		},
	})
	list = append(list, School{
		Name: "深圳",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:    30,
			Leader: "C",
		},
	})
	list = append(list, School{
		Name: "梅州",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:    5,
			Leader: "B",
		},
	})

	var newList []interface{}
	for _, v := range list {
		newList = append(newList, v)
	}

	type args struct {
		list      []interface{}
		condition *SortCondition
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "按照age从大到小",
			args: args{
				list: newList,
				condition: &SortCondition{
					Key:             "age",
					SortType:        SortDesc,
					KeyWeightGetter: keyWeightGetter,
				},
			},
			want: nil,
		},
		{
			name: "按照age从小到大",
			args: args{
				list: newList,
				condition: &SortCondition{
					Key:             "age",
					SortType:        SortAsc,
					KeyWeightGetter: keyWeightGetter,
				},
			},
			want: nil,
		},
		{
			name: "按照Leader从小到大ABC",
			args: args{
				list: newList,
				condition: &SortCondition{
					Key:             "leader",
					SortType:        SortAsc,
					KeyWeightGetter: keyWeightGetter,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginatorToolImplement{}
			got := p.SortList(tt.args.list, tt.args.condition)
			for _, v := range got {
				t.Logf("%+v", v)
			}
		})
	}
}

func TestPaginatorToolImplement_FilterList(t *testing.T) {

	var list []School
	list = append(list, School{
		Name: "广州",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:     18,
			Leader:  "A",
			Teacher: []string{"ming", "john"},
		},
	})
	list = append(list, School{
		Name: "潮州",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:     18,
			Leader:  "B",
			Teacher: []string{"ming", "john", "rose"},
		},
	})
	list = append(list, School{
		Name: "深圳",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:     30,
			Leader:  "C",
			Teacher: []string{"ming", "john", "jack"},
		},
	})
	list = append(list, School{
		Name: "梅州",
		Info: struct {
			Age     int
			Leader  string
			Teacher []string
		}{
			Age:     5,
			Leader:  "B",
			Teacher: []string{"lin"},
		},
	})

	var newList []interface{}
	for _, v := range list {
		newList = append(newList, v)
	}

	var keyValueGetter = func(listItem interface{}, key string) interface{} {
		if s, ok := listItem.(School); ok {
			switch key {
			case "age":
				return s.Info.Age
			case "leader":
				return s.Info.Leader
			case "teacher":
				return s.Info.Teacher
			}
		}
		return ""
	}

	var filterCondition []FilterCondition
	filterCondition = append(filterCondition, FilterCondition{
		Key:            "age",
		Operator:       Equals,
		Value:          18,
		KeyValueGetter: keyValueGetter,
	})
	filterCondition = append(filterCondition, FilterCondition{
		Key:            "leader",
		Operator:       Equals,
		Value:          "B",
		KeyValueGetter: keyValueGetter,
	})

	var filterConditionTeacher []FilterCondition
	filterConditionTeacher = append(filterConditionTeacher, FilterCondition{
		Key:            "teacher",
		Operator:       In,
		Value:          "jack",
		KeyValueGetter: keyValueGetter,
	})

	var filterConditionTeacherRose []FilterCondition
	filterConditionTeacherRose = append(filterConditionTeacherRose, FilterCondition{
		Key:            "teacher",
		Operator:       In,
		Value:          "rose",
		KeyValueGetter: keyValueGetter,
	})
	type args struct {
		list             []interface{}
		filterConditions []FilterCondition
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "测试获取age=18和leader=B的值",
			args: args{
				list:             newList,
				filterConditions: filterCondition,
			},
			want: nil,
		},
		{
			name: "测试获取teacher之中有jack的值",
			args: args{
				list:             newList,
				filterConditions: filterConditionTeacher,
			},
			want: nil,
		},
		{
			name: "测试获取teacher之中有rose的值",
			args: args{
				list:             newList,
				filterConditions: filterConditionTeacherRose,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginatorToolImplement{}
			got := p.FilterList(tt.args.list, tt.args.filterConditions)
			for _, v := range got {
				t.Logf("%+v", v)
			}
		})
	}
}

func TestPaginatorToolImplement_GetPage(t *testing.T) {
	list := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	type args struct {
		totalList []interface{}
		page      uint
		limit     uint
	}
	tests := []struct {
		name string
		args args
		want Paginator
	}{
		{
			name: "获取第1页，每页12条",
			args: args{
				totalList: list,
				page:      1,
				limit:     12,
			},
			want: Paginator{},
		},
		{
			name: "获取第10页，每页10条",
			args: args{
				totalList: list,
				page:      10,
				limit:     10,
			},
			want: Paginator{},
		},
		{
			name: "获取第2页，每页5条",
			args: args{
				totalList: list,
				page:      2,
				limit:     5,
			},
			want: Paginator{},
		},
		{
			name: "获取第3页，每页3条",
			args: args{
				totalList: list,
				page:      3,
				limit:     3,
			},
			want: Paginator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginatorToolImplement{}
			t.Logf("p=%+v", p.GetPage(tt.args.totalList, tt.args.page, tt.args.limit))
		})
	}
}
