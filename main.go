package main

import (
	"errors"
	"flag"
	"fmt"
	"markdown-parser/types"
	"os"
	"sort"
	"strconv"
	"strings"
)

var multilineTags = map[string]string{
	"```": "code",
}

var startTags = map[string]string{
	"---": "hr",
	"": "br",
}


var startWithTextTags = map[string]string{
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

func wrapStartTag(tag string, classes ...string) string{
	if len(classes) == 0{
		return "<"+tag+">"
	} else{
		htmlTag := "<"+tag+" class=\""
		for _, class := range classes{
			htmlTag += class+" "
		}
		htmlTag += "\">"
		return htmlTag
	}
}
func wrapEndTag(tag string) string{
	return "</"+tag+">"
}
func wrapWholeTag(tag string, text string) string{
	return wrapStartTag(tag)+text+wrapEndTag(tag)
}

func parse(markdownText string) (string, error){
	return  parseMultiLineTags(markdownText)	
}

func parseMultiLineTags(text string) (string, error){
	lines := strings.Split(text, "\n")
	stack := types.NewStack[string]()

	parsedText := ""

	for i ,line := range lines{
		justDidSmth := false
		for key := range multilineTags{
			if len(line) >= len(key) && key == line[:len(key)]{
				topElement := stack.TopElement();
				if topElement != nil && *topElement == key{
					parsedText+=wrapEndTag(multilineTags[key])+"\n"
					justDidSmth = true
					stack.Pop()
				}	else{
					parsedText+=wrapStartTag(multilineTags[key], line[len(key):])+"\n"
					justDidSmth = true
					stack.Push(key)
				}
				break
			}
		}
		if justDidSmth{
		}else if stack.TopElement() == nil{
			s, err := parseForStartTags(line, true);
			if err != nil{
				return "", errors.New(strconv.Itoa(i+1)+" line: "+err.Error())
			} 
			parsedText += s + "\n"
		} else{
			parsedText += line + "<br>"
		}
	}
	
	return parsedText, nil
}

func parseForStartTags(line string, wholeLine bool) (string, error){

	words := strings.SplitN(line," ", 2)	

	if tag, ok := startTags[line]; ok{
		return wrapStartTag(tag), nil
	} else if tag, ok := startWithTextTags[words[0]]; ok{
		parsed, err := parseForStartTags(words[1], false)
		if len(words) >= 2{
			return wrapWholeTag(tag, parsed), err
		}
		return "", errors.New("start tag "+tag+" requiers text")
		
	} else if wholeLine{
		parsed, err := parseForStartTags(line, false)
		return wrapWholeTag("p", parsed), err
	}
	// return line, nil
	return parseForDecor(line)
}

func tagListToTagSum(tags []string,start bool) string{
	tagSum := ""

	if start{
		for _, tag := range tags{
			tagSum+=wrapStartTag(tag)
		}
	} else{
		for i := len(tags)-1; i >= 0; i--{
			tagSum+=wrapEndTag(tags[i])
		}
	}
	
	return tagSum
}

func mapToKeyList[T any](obj map[string]T) []string{
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
					newText += tagListToTagSum(values, true)
				} else{
					stack.Pop()
					values := decorTags[key]
					newText += tagListToTagSum(values, false)
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


	parsed, err := parseMultiLineTags(s)

	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}else{
		fmt.Print(parsed)
	}
}


