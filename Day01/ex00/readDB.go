package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

type Recipe struct {
	Name        string `xml:"name" json:"name"`
	StoveTime   string `xml:"stovetime" json:"time"`
	Ingredients []struct {
		Name  string `xml:"itemname" json:"ingredient_name"`
		Count string `xml:"itemcount" json:"ingredient_count"`
		Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
	} `xml:"ingredients>item" json:"ingredients"`
}

type DBReader interface {
	readDB(filename string) ([]Recipe, error)
}

type XMLReader struct{}
type JSONReader struct{}

func (x *XMLReader) readDB(filename string) ([]Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var recipes struct {
		Recipes []Recipe `xml:"cake"`
	}
	err = xml.Unmarshal(data, &recipes)
	if err != nil {
		return nil, err
	}

	return recipes.Recipes, nil
}

func (j *JSONReader) readDB(filename string) ([]Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var recipes struct {
		Recipes []Recipe `json:"cake"`
	}

	err = json.Unmarshal(data, &recipes)
	if err != nil {
		return nil, err
	}

	return recipes.Recipes, nil
}

func getFilenameFromFlag() (*string, error) {
	filename := flag.String("f", "", "Имя файла базы данных (с расширением .xml или .json)")
	flag.Parse()

	if *filename == "" {
		return nil, errors.New("ERROR | Небходимо указать имя файла | Пример: go run readDB.go -f FILENAME.(json/lxm)")
	}

	return filename, nil
}

func getReader(filename *string) (DBReader, error) {
	var reader DBReader

	if strings.HasSuffix(*filename, ".xml") {
		reader = &XMLReader{}
	} else if strings.HasSuffix(*filename, ".json") {
		reader = &JSONReader{}
	} else {
		return nil, errors.New("ERROR | Неподдерживаемый формат файла")
	}

	return reader, nil
}

func printToConsole(recipes []Recipe, reader DBReader) error {
	if _, ok := reader.(*XMLReader); ok {
		jsonData, err := json.MarshalIndent(recipes, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	} else {
		xmlData, err := xml.MarshalIndent(recipes, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(xml.Header + string(xmlData))
	}
	return nil
}

func RunApplication() {
	filename, err := getFilenameFromFlag()
	if err != nil {
		fmt.Println(err)
		return
	}

	reader, err := getReader(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipes, err := reader.readDB(*filename)
	if err != nil {
		fmt.Println("ERROR | Ошибка при считывании файла:", err)
		return
	}

	err = printToConsole(recipes, reader)
	if err != nil {
		fmt.Println("ERROR | Ошибка при конвертации:", err)
		return
	}
}

func main() {
	RunApplication()
}
