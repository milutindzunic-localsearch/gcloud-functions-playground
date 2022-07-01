package onlim

import (
	"reflect"
	"testing"
)

// Note: the LocalEntries are incomplete, real ones contain much more data
var businessLocalEntry LocalEntry = []byte(`
		{
			"_id" : "X9CL4UiS2I6RqxwiLP75Tg",
			"addresses" : [
				{
					"address_types": [ "business", "restaurant" ],
					"business": {
						"categories": [
							{
								"id" : "75zEg3pk2lD7TdUrSzVkpQ"
							},
							{
								"id" : "_jyVPD-o3FF916UGAMIGsg"
							}
						],
						"secondary_categories": [
							{
								"id" : "u37d17CNjz9p2KDXMaPrLw"
							}
						]
					}
				}
			]
		}
	`)
var privateLocalEntry LocalEntry = []byte(`
		{
			"_id" : "3wYUIlZ1gWeHhwShKW050g",
			"addresses" : [
				{
					"address_types": [ "private" ]
				}
			],
		}
	`)

func TestLocalEntry_IsBusiness_ReturnsTrueForBusinessEntry(t *testing.T) {
	isBusiness := businessLocalEntry.IsBusiness()

	if !isBusiness {
		t.Error("IsBusiness() returns false for a business LocalEntry")
	}
}

func TestLocalEntry_IsBusiness_ReturnsFalseForPrivateEntry(t *testing.T) {
	isBusiness := privateLocalEntry.IsBusiness()

	if isBusiness {
		t.Error("IsBusiness() returns true for a private LocalEntry")
	}
}

func TestLocalEntry_HasOneOf_ReturnsTrueWhenAnyCategoryMatches(t *testing.T) {
	catIDs := []CategoryID{"_jyVPD-o3FF916UGAMIGsg"}

	if !businessLocalEntry.HasOneOf(catIDs) {
		t.Error("HasOneOf() returns false even though there is a category match")
	}
}

func TestLocalEntry_HasOneOf_ReturnsFalseWhenNoCategoriesMatch(t *testing.T) {
	catIDs := []CategoryID{"foobar"}

	if businessLocalEntry.HasOneOf(catIDs) {
		t.Error("HasOneOf() returns true even though there are no category matches")
	}
}

func TestLocalEntry_primaryCategoryIDs(t *testing.T) {
	expectedCatIDs := []CategoryID{"75zEg3pk2lD7TdUrSzVkpQ", "_jyVPD-o3FF916UGAMIGsg"}
	catIDs := businessLocalEntry.primaryCategoryIDs()

	if !reflect.DeepEqual(catIDs, expectedCatIDs) {
		t.Errorf("primaryCategoryIDs() was incorrect, got: %s, want: %s.", catIDs, expectedCatIDs)
	}
}

func TestLocalEntry_secondaryCategoryIDs(t *testing.T) {
	expectedCatIDs := []CategoryID{"u37d17CNjz9p2KDXMaPrLw"}
	catIDs := businessLocalEntry.secondaryCategoryIDs()

	if !reflect.DeepEqual(catIDs, expectedCatIDs) {
		t.Errorf("secondaryCategoryIDs() was incorrect, got: %s, want: %s.", catIDs, expectedCatIDs)
	}
}
