package recipes

import (
	"encoding/json"
	"github.com/dmitriibb/go-common/restaurant-common/model"
	"net/http"
)

func GetAllRecipesAsMenu(w http.ResponseWriter, _ *http.Request) {
	allRecipes, err := GetAllRecipes()
	if err != nil {
		logger.Error("can't get all recipies because '%v'", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		resp := model.CommonErrorResponse{
			Type:    model.CommonErrorTypeInvalidData,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	menuItems := make([]*model.MenuItemDto, 0, len(allRecipes))
	for _, recipe := range allRecipes {
		menuItems = append(menuItems,
			&model.MenuItemDto{
				Name:        recipe.Name,
				Description: recipe.Description,
				Ingredients: recipe.Ingredients,
			},
		)
	}
	json.NewEncoder(w).Encode(model.MenuDto{menuItems})
}
