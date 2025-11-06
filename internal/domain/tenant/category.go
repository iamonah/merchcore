package tenant

import (
	"fmt"
	"strings"
)

type BusinessCategory string

var (
	CategoryFashion       = newBusinessCategory("fashion")
	CategoryElectronics   = newBusinessCategory("electronics")
	CategoryBeauty        = newBusinessCategory("beauty")
	CategoryFood          = newBusinessCategory("food")
	CategoryHome          = newBusinessCategory("home")
	CategoryHealth        = newBusinessCategory("health")
	CategoryBooks         = newBusinessCategory("books")
	CategoryAutomotive    = newBusinessCategory("automotive")
	CategoryToys          = newBusinessCategory("toys")
	CategoryArt           = newBusinessCategory("art")
	CategoryJewelry       = newBusinessCategory("jewelry")
	CategoryServices      = newBusinessCategory("services")
	CategoryPets          = newBusinessCategory("pets")
	CategoryEducation     = newBusinessCategory("education")
	CategoryAgriculture   = newBusinessCategory("agriculture")
	CategoryRealEstate    = newBusinessCategory("real_estate")
	CategoryEntertainment = newBusinessCategory("entertainment")
	CategoryTechnology    = newBusinessCategory("technology")
)

var categories = make(map[string]BusinessCategory)

func newBusinessCategory(v string) BusinessCategory {
	bc := BusinessCategory(v)
	categories[strings.ToLower(v)] = bc
	return bc
}

func ParseBusinessCategory(category string) (BusinessCategory, error) {
	v, ok := categories[category]
	if !ok {
		return "", fmt.Errorf("invalid businesscategory: %v", category)
	}
	return v, nil
}

func AllBusinessCategories() []BusinessCategory {
	out := make([]BusinessCategory, 0, len(categories))
	for _, c := range categories {
		out = append(out, c)
	}
	return out
}
