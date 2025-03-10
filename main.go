package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type SchemaBuilder struct {
	writer io.Writer
	indent string
}

func NewSchemaBuilder(w io.Writer) *SchemaBuilder {
	return &SchemaBuilder{
		writer: w,
		indent: "  ",
	}
}

func (sb *SchemaBuilder) buildSchema(v interface{}, level int) error {
	if v == nil {
		_, err := fmt.Fprint(sb.writer, "null")
		return err
	}

	switch val := v.(type) {
	case map[string]interface{}:
		return sb.handleObject(val, level)
	case []interface{}:
		return sb.handleArray(val, level)
	case string:
		return sb.writeType("string")
	case float64:
		return sb.writeType("number")
	case bool:
		return sb.writeType("boolean")
	default:
		return sb.writeType(fmt.Sprintf("unknown (%T)", v))
	}
}

func (sb *SchemaBuilder) handleObject(m map[string]interface{}, level int) error {
	if _, err := fmt.Fprintln(sb.writer, "{"); err != nil {
		return err
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	for i, key := range keys {
		indent := strings.Repeat(sb.indent, level+1)
		if _, err := fmt.Fprintf(sb.writer, "%s%q: ", indent, key); err != nil {
			return err
		}

		if err := sb.buildSchema(m[key], level+1); err != nil {
			return err
		}

		if i < len(keys)-1 {
			if _, err := fmt.Fprint(sb.writer, ","); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(sb.writer); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(sb.writer, "%s}", strings.Repeat(sb.indent, level)); err != nil {
		return err
	}
	return nil
}

func (sb *SchemaBuilder) handleArray(arr []interface{}, level int) error {
	if _, err := fmt.Fprintln(sb.writer, "["); err != nil {
		return err
	}

	if len(arr) > 0 {
		indent := strings.Repeat(sb.indent, level+1)
		if _, err := fmt.Fprint(sb.writer, indent); err != nil {
			return err
		}
		if err := sb.buildSchema(arr[0], level+1); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(sb.writer); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(sb.writer, "%s]", strings.Repeat(sb.indent, level)); err != nil {
		return err
	}
	return nil
}

func (sb *SchemaBuilder) writeType(t string) error {
	_, err := fmt.Fprint(sb.writer, t)
	return err
}

func processInput(input io.Reader, output io.Writer) error {
	var data interface{}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	builder := NewSchemaBuilder(output)
	if err := builder.buildSchema(data, 0); err != nil {
		return fmt.Errorf("building schema: %w", err)
	}

	if _, err := fmt.Fprintln(output); err != nil {
		return fmt.Errorf("writing final newline: %w", err)
	}

	return nil
}

func main() {
	input := os.Stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	if err := processInput(input, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
