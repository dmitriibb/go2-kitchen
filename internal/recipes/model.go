package recipes

// RecipeStage Used to store recipes and for tracing dish cooking progress
type RecipeStage struct {
	Name          string            `json:"name" bson:"name"`
	Description   string            `json:"description" bson:"description"`
	Ingredients   []string          `json:"ingredients" bson:"ingredients"`
	TimeToWaitSec int64             `json:"timeToWaitSec" bson:"timeToWaitSec"`
	TimeStarted   int64             `json:"timeStarted" bson:"timeStarted"`
	Status        RecipeStageStatus `json:"status" bson:"status"`
	Comment       string            `json:"comment" bson:"comment"`
	SubStages     []*RecipeStage    `json:"subStages" bson:"subStages"`
}

type RecipeStageStatus string

const (
	RecipeStageStatusEmpty      RecipeStageStatus = ""
	RecipeStageStatusInProgress RecipeStageStatus = "InProgress"
	RecipeStageStatusFinished   RecipeStageStatus = "Finished"
	RecipeStageStatusError      RecipeStageStatus = "Error"
)
