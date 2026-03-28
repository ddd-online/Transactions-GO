package operator

import "github.com/billadm/models/dto"

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type SortField struct {
	Field string
	Order SortOrder
}

// sortableTrDtos 实现 sort.Interface，并携带排序规则
type sortableTrDtos struct {
	data       []*dto.TransactionRecordDto
	sortFields []SortField
}

func (s sortableTrDtos) Len() int {
	return len(s.data)
}

func (s sortableTrDtos) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s sortableTrDtos) Less(i, j int) bool {
	a := s.data[i]
	b := s.data[j]

	for _, sf := range s.sortFields {
		switch sf.Field {
		case "price":
			if a.Price != b.Price {
				if sf.Order == Asc {
					return a.Price < b.Price
				}
				return a.Price > b.Price
			}
		case "transactionType":
			if a.TransactionType != b.TransactionType {
				if sf.Order == Asc {
					return a.TransactionType < b.TransactionType
				}
				return a.TransactionType > b.TransactionType
			}
		case "category":
			if a.Category != b.Category {
				if sf.Order == Asc {
					return a.Category < b.Category
				}
				return a.Category > b.Category
			}
		case "transactionAt":
			if a.TransactionAt != b.TransactionAt {
				if sf.Order == Asc {
					return a.TransactionAt < b.TransactionAt
				}
				return a.TransactionAt > b.TransactionAt
			}
		default:
			continue // 忽略未知字段
		}
	}
	return false
}
