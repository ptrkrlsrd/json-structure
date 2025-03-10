package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Indent struct {
	Level      int
	IndentChar rune
}

func DefaultIndent() Indent {
	return Indent{
		Level:      0,
		IndentChar: ' ',
	}
}

func (i Indent) addIndent() Indent {
	i.Level = +1

	return i
}

func main() {
	var input io.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	var jsonData interface{}
	decoder := json.NewDecoder(input)
	err := decoder.Decode(&jsonData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	var sb strings.Builder
	buildStructure(&sb, jsonData, 0)
	sb.WriteString("\n")
	fmt.Print(sb.String())
}

func buildStructure(sb *strings.Builder, v interface{}, indent int) {
	if v == nil {
		sb.WriteString("null")
		return
	}

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Map:
		sb.WriteString("{\n")
		m := reflect.ValueOf(v)
		keys := m.MapKeys()
		for i, key := range keys {
			sb.WriteString(fmt.Sprintf("%s%v: ", strings.Repeat("  ", indent+1), key.Interface()))
			buildStructure(sb, m.MapIndex(key).Interface(), indent+1)
			if i < len(keys)-1 {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}
		sb.WriteString(strings.Repeat("  ", indent) + "}")
	case reflect.Slice:
		sb.WriteString("[\n")
		s := reflect.ValueOf(v)
		if s.Len() > 0 {
			sb.WriteString(strings.Repeat("  ", indent+1))
			buildStructure(sb, s.Index(0).Interface(), indent+1)
			sb.WriteString("\n")
		}
		sb.WriteString(strings.Repeat("  ", indent) + "]")
	case reflect.String:
		sb.WriteString("string")
	case reflect.Float64:
		sb.WriteString("number")
	case reflect.Bool:
		sb.WriteString("boolean")
	case reflect.Invalid:
		sb.WriteString("null")
	default:
		sb.WriteString(fmt.Sprintf("unknown (%v)", t.Kind()))
	}
}
