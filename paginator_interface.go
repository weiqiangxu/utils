package utils

type Operator string

const (
	NotEquals    Operator = "!="
	Equals       Operator = "="
	In           Operator = "in"
	NotIn        Operator = "notin"  // TODO implement
	DoesNotExist Operator = "!"      // TODO implement
	DoubleEquals Operator = "=="     // TODO implement
	Exists       Operator = "exists" // TODO implement
	GreaterThan  Operator = "gt"     // TODO implement
	LessThan     Operator = "lt"     // TODO implement
)

type SortType string

const (
	SortDesc SortType = "desc"
	SortAsc  SortType = "asc"
)

type SortCondition struct {
	Key             string
	SortType        SortType
	KeyWeightGetter KeyWeightGetterFunc
}

type FilterCondition struct {
	Key            string
	Operator       Operator
	Value          interface{}
	KeyValueGetter KeyValueGetterFunc
}

type Paginator struct {
	Page     uint          `json:"page" form:"page"`         // 页码
	PageSize uint          `json:"pageSize" form:"pageSize"` // 每页大小
	Total    uint          `json:"total" form:"total"`       // 总条数
	List     []interface{} `json:"list" form:"list"`         // 页数据
}

// KeyValueGetterFunc 键值获取
type KeyValueGetterFunc func(listItem interface{}, key string) interface{}

// KeyWeightGetterFunc 键的权重获取
type KeyWeightGetterFunc func(listItem interface{}, key string) int

type Tool interface {
	// SortList 按字段排序
	SortList(totalList []interface{}, condition *SortCondition) []interface{}
	// FilterList 过滤特定字段值
	FilterList(totalList []interface{}, filterConditions []FilterCondition) []interface{}
	// GetPage 对list分页并获取指定页数
	GetPage(totalList []interface{}, page uint, limit uint) Paginator
	// GetFilterPage 对list过滤并排序后获取指定页的数据
	GetFilterPage(totalList []interface{}, page uint, limit uint, filterConditions []FilterCondition, sortCondition *SortCondition) Paginator
}
