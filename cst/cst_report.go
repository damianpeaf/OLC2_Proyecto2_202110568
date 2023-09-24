package cst

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type CSTresponse struct {
	SVGTree string `json:"svgtree"`
}

func ReadFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, _ := io.ReadAll(file)
	return string(content)
}

func CstReport(input string) string {

	// get the content (relative to this path) ../compiler/TSwiftLanguage.g4

	parserContent := ""

	_, filename, _, _ := runtime.Caller(0)

	path := filepath.Dir(filename)

	// remove \cst from the path
	path = path[:len(path)-4]

	parser, err := json.Marshal(ReadFile(filepath.Join(path, "/compiler/TSwiftLanguage.g4")))

	if err != nil {
		fmt.Println(err)
	}
	parserContent = string(parser)

	lexerContent := ""
	lexer, err := json.Marshal(ReadFile(filepath.Join(path, "/compiler/TSwiftLexer.g4")))

	if err != nil {
		fmt.Println(err)
	}

	lexerContent = string(lexer)

	jinput, err := json.Marshal(input)
	finput := string(jinput)

	payload := []byte(
		fmt.Sprintf(
			`{
				"grammar": %s,
				"input": %s,
				"lexgrammar": %s,
				"start": "%s"
			}`,
			parserContent,
			finput,
			lexerContent,
			"program",
		),
	)

	req, err := http.NewRequest("POST", "http://lab.antlr.org/parse/", bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return ""
	}

	// create a map to store the json
	var data map[string]interface{}

	// // unmarshal the json
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling json:", err)
		return ""
	}

	result := data["result"].(map[string]interface{})

	return result["svgtree"].(string)
}
