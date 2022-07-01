package onlim

import (
	"fmt"
	"github.com/spyzhov/ajson"
)

// LocalEntry - we care only for a few properties -
// so we store it as unmarshalled json, and use jsonpath to drill into what we need.
type LocalEntry []byte

type CategoryID string

const addressTypeBusiness = "business"

func (entry LocalEntry) IsBusiness() bool {
	return entry.hasAddressType(addressTypeBusiness)
}

func (entry LocalEntry) HasOneOf(categoryIDs []CategoryID) bool {
	cIDs := append(entry.primaryCategoryIDs(), entry.secondaryCategoryIDs()...)

	for _, categoryID := range categoryIDs {
		for _, cID := range cIDs {
			if categoryID == cID {
				return true
			}
		}
	}

	return false
}

func (entry LocalEntry) hasAddressType(addressType string) bool {
	addressTypeJsonPath := fmt.Sprintf("$.addresses[0].address_types[?(@ == '%s')]", addressType)

	nodes, err := ajson.JSONPath(entry, addressTypeJsonPath)

	return err == nil && len(nodes) > 0
}

func (entry LocalEntry) primaryCategoryIDs() []CategoryID {
	categoryIDsJsonPath := fmt.Sprintf(`$.addresses[0].business.categories[*].id`)
	return entry.getCategoryIDs(categoryIDsJsonPath)
}

func (entry LocalEntry) secondaryCategoryIDs() []CategoryID {
	secondaryCategoryIDsJsonPath := fmt.Sprintf(`$.addresses[0].business.secondary_categories[*].id`)
	return entry.getCategoryIDs(secondaryCategoryIDsJsonPath)
}

func (entry LocalEntry) getCategoryIDs(jsonPath string) []CategoryID {

	nodes, err := ajson.JSONPath(entry, jsonPath)
	if err != nil || len(nodes) <= 0 {
		return []CategoryID{}
	}

	var categoryIDs []CategoryID
	for _, node := range nodes {
		cid, err := node.GetString()
		if err == nil {
			categoryIDs = append(categoryIDs, CategoryID(cid))
		}
	}

	return categoryIDs
}
