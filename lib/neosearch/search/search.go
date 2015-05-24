package search

import (
	"errors"
	"fmt"

	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
)

type (
	DSL       map[string]interface{}
	Result    string
	ResultSet []Result
)

func (d DSL) Map() map[string]interface{} { return map[string]interface{}(d) }

func Search(ind *index.Index, dsl DSL, limit uint) ([]string, uint64, error) {
	var (
		listOp        []interface{}
		hasAnd, hasOr bool
		resultDocIDs  []uint64
	)

	listOp, hasAnd = dsl["$and"].([]interface{})

	if !hasAnd {
		listOp, hasOr = dsl["$or"].([]interface{})
	}

	if !hasAnd && !hasOr {
		return nil, 0, errors.New("Invalid search DSL. No $and or $or clause found.")
	}

	for _, clause := range listOp {
		filter, ok := clause.(map[string]interface{})

		if !ok {
			return nil, 0, fmt.Errorf("Invalid clause '%s'.", clause)
		}

		field, value := getFieldValue(filter)

		if field == "" || value == nil {
			return nil, 0, fmt.Errorf("Invalid clause '%s'.", clause)
		}

		strValue, ok := value.(string)

		if !ok {
			return nil, 0, fmt.Errorf("Invalid field value: %s", value)
		}

		docIDs, _, err := ind.FilterTermID([]byte(field), []byte(strValue), 0)

		if err != nil {
			return nil, 0, err
		}

		if len(resultDocIDs) == 0 {
			resultDocIDs = docIDs
			continue
		}

		if hasAnd {
			resultDocIDs = and(resultDocIDs, docIDs)
		}
	}

	results, err := ind.GetDocs(resultDocIDs, limit)
	return results, uint64(len(resultDocIDs)), err
}

// TODO: we need benchmark this algorithm and optimize
func and(a, b []uint64) []uint64 {
	var (
		aSize, bSize          = len(a), len(b)
		maxSize, i, j, resIdx = 0, 0, 0, 0
		result                []uint64
	)

	if aSize > bSize {
		maxSize = aSize
	} else {
		maxSize = bSize
	}

	result = make([]uint64, maxSize)

	for i < aSize && j < bSize {
		if a[i] == b[j] {
			result[resIdx] = a[i]
			i++
			j++
			resIdx++
		} else if a[i] < b[j] {
			i++
			continue
		} else if a[i] > b[j] {
			j++
			continue
		}
	}

	return result[0:resIdx]
}

func getFieldValue(filter map[string]interface{}) (string, interface{}) {
	for field, value := range filter {
		return field, value
	}

	return "", nil
}
