package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func readSnapshotAndFillHashmap(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hashmap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		hash := generateMD5Hash(path)
		hashmap[hash] = path
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashmap, nil
}

func compareSnapshots(snapshot1, snapshot2 string) {
	hashmap, err := readSnapshotAndFillHashmap(snapshot1)
	if err != nil {
		fmt.Println("Ошибка при чтении snapshot1:", err)
		return
	}

	file, err := os.Open(snapshot2)
	if err != nil {
		fmt.Println("Ошибка при чтении snapshot2:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		hash := generateMD5Hash(path)
		if _, exists := hashmap[hash]; exists {
			delete(hashmap, hash)
		} else {
			fmt.Println("ADDED", path)
		}
	}

	for _, path := range hashmap {
		fmt.Println("REMOVED", path)
	}
}

func RunApplication() {
	oldSnapshot := flag.String("old", "", "Путь к старому снимку файловой системы")
	newSnapshot := flag.String("new", "", "Путь к новому снимку файловой системы")

	flag.Parse()

	if *oldSnapshot == "" || *newSnapshot == "" {
		fmt.Println("ERROR | Необходимо указать пути к обоим снимкам файловой системы.")
		return
	}

	compareSnapshots(*oldSnapshot, *newSnapshot)
}

func main() {
	RunApplication()
}
