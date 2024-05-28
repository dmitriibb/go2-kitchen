package recipes

import (
	"context"
	"encoding/json"
	"fmt"
	commonMongo "github.com/dmitriibb/go-common/db/mongo"
	"github.com/dmitriibb/go-common/logging"
	commonInitializer "github.com/dmitriibb/go-common/utils/initializer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	recipesCollection = "recipes"
	saveModeInMemory  = "inMemory"
	saveModeMongo     = "mongo"
)

var logger = logging.NewLogger("RecipesRepository")

var saveMode = saveModeMongo

// store in json
var inMemoryData = make(map[string]string)
var initializer = commonInitializer.New(logger)
var dbName = ""

func Init() {
	initializer.Init(func() error {
		dbName = commonMongo.GetDbName()
		initData()
		return nil
	})
}

func initData() {
	if saveMode == saveModeMongo {
		createMongoCollection()
	}
	save(recipeBurger())
	save(recipeCola())
	save(recipeCoffee())
	save(recipeWater())
	save(recipeBread())
	save(recipePasta())
}

func save(recipe RecipeStage) {
	if saveMode == saveModeInMemory {
		saveInMemory(recipe)
	} else {
		saveInMongo(recipe)
	}
}

func saveInMemory(recipe RecipeStage) {
	v, _ := json.Marshal(recipe)
	inMemoryData[recipe.Name] = string(v)
}

func createMongoCollection() {
	client := commonMongo.GetClient()
	defer client.Disconnect(context.TODO())
	err := client.Database(dbName).CreateCollection(context.TODO(), recipesCollection)
	if err != nil {
		logger.Warn("Collection '%s' created with error %s", recipesCollection, err.Error())
	} else {
		logger.Info("Collection '%s' created", recipesCollection)
	}
}

func saveInMongo(recipe RecipeStage) {
	client := commonMongo.GetClient()
	defer client.Disconnect(context.TODO())
	collection := client.Database(dbName).Collection(recipesCollection)

	filter := bson.D{{"name", recipe.Name}}
	update := bson.D{{"$set", recipe}}
	result, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	//result, err := collection.InsertOne(context.TODO(), recipe)
	if err != nil {
		logger.Error("can't save recipe in DB because %v", err.Error())
		return
	}
	logger.Info("saved recipe '%v' in DB. Result: %v", recipe.Name, result.UpsertedID)
}

func GetRecipe(name string) (RecipeStage, error) {
	if saveMode == saveModeInMemory {
		res, err := getRecipeFromMemory(name)
		return *res, err
	} else {
		res, err := getRecipeFromMongo(name)
		return *res, err
	}
}

func GetAllRecipes() ([]*RecipeStage, error) {
	if saveMode == saveModeInMemory {
		panic("not implemented an never will be")
	} else {
		res, err := getAllFromMongo()
		return res, err
	}
	//return nil, errors.New("Fake error")
}

func getRecipeFromMemory(name string) (*RecipeStage, error) {
	data, ok := inMemoryData[name]
	if !ok {
		return &RecipeStage{}, fmt.Errorf("recipe for '%v' not found", name)
	}
	var recipe RecipeStage
	err := json.Unmarshal([]byte(data), &recipe)
	if err != nil {
		return &RecipeStage{}, err
	}
	return &recipe, nil
}

func getRecipeFromMongo(name string) (*RecipeStage, error) {
	ctx := context.TODO()
	filter := bson.D{{"name", name}}
	var err error
	var cur *mongo.Cursor
	var res *RecipeStage
	f := func(client *mongo.Client) any {
		collection := client.Database(commonMongo.GetDbName()).Collection(recipesCollection)
		cur, err = collection.Find(ctx, filter)
		if err != nil {
			logger.Error("can't find '%v' in the %v", name, recipesCollection)
			return nil
		}
		var results []RecipeStage
		if err = cur.All(ctx, &results); err != nil {
			logger.Error("can't parse cursor into []RecipeStage")
			return nil
		}
		if len(results) > 1 {
			logger.Warn("found %v recipes for '%v'", len(results), name)
			res = &results[0]
		} else if len(results) == 1 {
			res = &results[0]
		} else {
			logger.Error("can't find recipe for %s", name)
		}
		return nil
	}
	commonMongo.UseClient(ctx, f)
	return res, err
}

func getAllFromMongo() ([]*RecipeStage, error) {
	ctx := context.TODO()
	var err error
	var cur *mongo.Cursor
	var results []RecipeStage
	f := func(client *mongo.Client) any {
		collection := client.Database(commonMongo.GetDbName()).Collection(recipesCollection)
		cur, err = collection.Find(ctx, bson.D{})
		if err != nil {
			logger.Error("can't find all'%v'", recipesCollection)
			return nil
		}
		if err = cur.All(ctx, &results); err != nil {
			logger.Error("can't parse cursor into []RecipeStage")
			return nil
		}
		return nil
	}
	commonMongo.UseClient(ctx, f)
	resultPointers := make([]*RecipeStage, 0)
	for _, rs := range results {
		resultPointers = append(resultPointers, &rs)
	}
	return resultPointers, err
}

func recipeBurger() RecipeStage {
	burger := RecipeStage{}
	burger.Name = "burger"
	burger.Ingredients = []string{}
	burger.SubStages = []*RecipeStage{
		{
			Name:        "cut vegetables",
			Ingredients: []string{"tomato", "lettuce", "onion"},
		},
		{
			Name:          "grill meet",
			Ingredients:   []string{"beef"},
			TimeToWaitSec: 20,
		},
		{
			Name:        "assemble burger",
			Ingredients: []string{"burger bun", "mayo", "cheese"},
		},
	}
	return burger
}
func recipeCola() RecipeStage {
	cola := RecipeStage{}
	cola.Name = "cola"
	cola.Ingredients = []string{"cola", "ice"}
	cola.SubStages = []*RecipeStage{
		{
			Name:        "open cola",
			Ingredients: []string{"cola"},
		},
		{
			Name:        "add ice",
			Ingredients: []string{"ice"},
		},
	}
	return cola
}
func recipeCoffee() RecipeStage {
	coffee := RecipeStage{}
	coffee.Name = "coffee"
	coffee.Ingredients = []string{"coffee", "milk"}
	coffee.SubStages = []*RecipeStage{
		{
			Name:        "brew coffee",
			Ingredients: []string{"coffee"},
		},
		{
			Name:        "add milk",
			Ingredients: []string{"milk"},
		},
	}
	return coffee
}
func recipeWater() RecipeStage {
	water := RecipeStage{}
	water.Name = "water"
	water.Ingredients = []string{"water"}
	return water
}
func recipeBread() RecipeStage {
	water := RecipeStage{}
	water.Name = "bread"
	water.Ingredients = []string{"bread"}
	return water
}
func recipePasta() RecipeStage {
	pasta := RecipeStage{}
	pasta.Name = "pasta"
	pasta.Ingredients = []string{"egg", "flour", "onion", "chicken", "salt", "butter"}
	pasta.SubStages = []*RecipeStage{
		{
			Name:        "make pasta",
			Ingredients: []string{"egg", "flour"},
			SubStages: []*RecipeStage{
				{
					Name:        "mix eggs and flour",
					Ingredients: []string{"egg", "flour"},
				},
				{
					Name:          "boil water",
					TimeToWaitSec: 10,
				},
				{
					Name:          "cook pasta in water",
					TimeToWaitSec: 12,
				},
			},
		},
		{
			Name:        "fry chicken",
			Ingredients: []string{"chicken", "onion"},
		},
		{
			Name:        "fry chicken with pasta",
			Ingredients: []string{"salt", "butter"},
		},
	}
	return pasta
}
