package main

import (
	"errors"
	"flag"
	"fmt"
	"markdown-parser/types"
	"os"
	"sort"
	"strings"
)

var startTags = map[string]string{
	"#": "h1",
	"##": "h2",
	"###": "h3",
	"####": "h4",
	"#####": "h5",
	"######": "h6",
}

var decorTags = map[string][]string{
	"*": {"em"},
	"-": {"em"},
	"**": {"strong"},
	"--": {"strong"},
	"~~": {"s"},
	"***": {"strong", "em"},
	"---": {"strong", "em"},
	"`" : {"code"},
}


func readFile(path string) (string,error){
	file, err := os.ReadFile(path)
    if err != nil {
			return "", err
    }

	s := string(file);


	return s, nil;
}

var path string;

func loadVariables() error{
	pathRef := flag.String("path", "", "path to markdown file")

	flag.Parse()
	
	if (*pathRef == ""){
		return errors.New("missing path variable")
	}

	path = *pathRef
	return nil;
}

func parse(line string, wholeLine bool) (string, error){

	words := strings.SplitN(line," ", 2)	

	if(line == ""){
		return "</br>", nil
	} else if value, ok := startTags[words[0]]; ok{
		parsed, err := parse(words[1], false)
		return "<"+value+">"+parsed+"</"+value+">", err
	} else if wholeLine{
		parsed, err := parse(line, false)
		return "<p>"+parsed+"</p>", err
	}
	// return line, nil
	return parseForDecor(line)
}

func valuesToTag(values []string,start bool) string{
	tags := ""

	if start{
		for _, element := range values{
			tags+="<"+element+">"
		}
	} else{
		for i := len(values)-1; i >= 0; i--{
			tags+="</"+values[i]+">"
		}
	}
	
	return tags
}

func mapToKeyList(obj map[string][]string) []string{
	var keyList []string
	for key := range decorTags{
		keyList = append(keyList, key)
	}

	sort.Slice(keyList, func(i, j int) bool {
		return len(keyList[i]) > len(keyList[j])
	})

	return keyList
}

func parseForDecor(text string) (string, error){
	stack := types.NewStack[string]()
	newText := ""
	chars := []rune(text)
	keyList := mapToKeyList(decorTags)
	for i := 0; i < len(chars); i++{
		key_found := false
		for _, key := range keyList{
			if(len(chars)-i >= len(key) && key == string(chars[i:i+len(key)])){
				key_found = true;
				i += len(key)-1
				topElement := stack.TopElement()
				if(topElement == nil || *topElement != key){
					stack.Push(key)
					values := decorTags[key]
					newText += valuesToTag(values, true)
				} else{
					stack.Pop()
					values := decorTags[key]
					newText += valuesToTag(values, false)
				}
				break
			} 	
		}	
		if !key_found{
			newText+=string(chars[i])
		}
	}

	if len(stack.Elements) != 0{
		return "", errors.New("decor aply error")
	}
	return newText, nil
}

func main(){
	err := loadVariables()

	if (err != nil){
		fmt.Println(err.Error())
		os.Exit(1)
	}

	
	s, err := readFile(path)

	if (err != nil){
		fmt.Println(err.Error())
	}

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		parsedString, err := parse(line, true)
		if (err != nil){
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(parsedString)
	}

}


