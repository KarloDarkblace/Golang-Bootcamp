package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Ingredient struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type Recipe struct {
	Name        string       `xml:"name" json:"name"`
	StoveTime   string       `xml:"stovetime" json:"time"`
	Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}

type DBReader interface {
	ReadDB(filename string) ([]Recipe, error)
}

type XMLReader struct{}
type JSONReader struct{}

func (x *XMLReader) ReadDB(filename string) ([]Recipe, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var recipes struct {
		Recipes []Recipe `xml:"cake"`
	}
	if err := xml.Unmarshal(data, &recipes); err != nil {
		return nil, err
	}

	return recipes.Recipes, nil
}

func (j *JSONReader) ReadDB(filename string) ([]Recipe, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var recipes struct {
		Recipes []Recipe `json:"cake"`
	}

	if err := json.Unmarshal(data, &recipes); err != nil {
		return nil, err
	}

	return recipes.Recipes, nil
}

func GetFilenameFromFlag() (string, string, error) {
	oldFile := flag.String("old", "", "Путь к оригинальному файлу базы данных (XML или JSON)")
	newFile := flag.String("new", "", "Путь к новому файлу базы данных (XML или JSON)")

	flag.Parse()

	if *oldFile == "" || *newFile == "" {
		return "", "", errors.New("ERROR | Необходимо указать пути к обоим файлам баз данных")
	}

	return *oldFile, *newFile, nil
}

func GetReader(filename string) (DBReader, error) {
	if strings.HasSuffix(filename, ".xml") {
		return &XMLReader{}, nil
	} else if strings.HasSuffix(filename, ".json") {
		return &JSONReader{}, nil
	}
	return nil, errors.New("ERROR | Неподдерживаемый формат файла")
}

func CompareDatabases(oldRecipes []Recipe, newRecipes []Recipe) {
	oldRecipesMap := map[string]Recipe{}
	newRecipesMap := map[string]Recipe{}

	for _, recipe := range oldRecipes {
		oldRecipesMap[recipe.Name] = recipe
	}
	for _, recipe := range newRecipes {
		newRecipesMap[recipe.Name] = recipe
	}

	for name, newRecipe := range newRecipesMap {
		if oldRecipe, exists := oldRecipesMap[name]; !exists {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		} else {
			CompareRecipes(oldRecipe, newRecipe)
		}
	}

	for name := range oldRecipesMap {
		if _, exists := newRecipesMap[name]; !exists {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
		}
	}
}

func CompareRecipes(oldRecipe Recipe, newRecipe Recipe) {
	if oldRecipe.StoveTime != newRecipe.StoveTime {
		fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", newRecipe.Name, newRecipe.StoveTime, oldRecipe.StoveTime)
	}

	CompareIngredients(oldRecipe, newRecipe)
}

func CompareIngredients(oldRecipe Recipe, newRecipe Recipe) {
	oldIngredientsMap := map[string]Ingredient{}
	for _, ingredient := range oldRecipe.Ingredients {
		oldIngredientsMap[ingredient.Name] = ingredient
	}

	for _, newIngredient := range newRecipe.Ingredients {
		oldIngredient, exists := oldIngredientsMap[newIngredient.Name]
		if !exists {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", newIngredient.Name, newRecipe.Name)
			continue
		}

		if oldIngredient.Count != newIngredient.Count {
			fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", newIngredient.Name, newRecipe.Name, newIngredient.Count, oldIngredient.Count)
		}
		if oldIngredient.Unit != newIngredient.Unit {
			fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", newIngredient.Name, newRecipe.Name, newIngredient.Unit, oldIngredient.Unit)
		}
		delete(oldIngredientsMap, newIngredient.Name)
	}

	for name := range oldIngredientsMap {
		fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", name, oldRecipe.Name)
	}
}

func RunApplication() {
	oldFilename, newFilename, err := GetFilenameFromFlag()
	if err != nil {
		fmt.Println(err)
		return
	}

	oldReader, err := GetReader(oldFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	newReader, err := GetReader(newFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	oldRecipes, err := oldReader.ReadDB(oldFilename)
	if err != nil {
		fmt.Println("ERROR | Ошибка при считывании файла:", err)
		return
	}

	newRecipes, err := newReader.ReadDB(newFilename)
	if err != nil {
		fmt.Println("ERROR | Ошибка при считывании файла:", err)
		return
	}

	CompareDatabases(oldRecipes, newRecipes)
}

func main() {
	RunApplication()
}
