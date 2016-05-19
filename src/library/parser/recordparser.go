package parser

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//tsv, csv parser
type Parser struct {
	FieldsSplitter string //field splitter in tsv or csv file: \t or , or etc
	LinesSplitter  string //file line splitter: \n
	PanicOnError   bool   //if true, panic
	Tag            string //tag name for data struct field, ex. header

	dataTypeParser map[reflect.Type]reflect.Value
	//key: data type TYPE,
	//value: parse function of proto func(string) (TYPE, error)
}

type IAfterParse interface {
	AfterParse()
}

//parse function of proto func(string) (TYPE, error)
func (p *Parser) RegisterParser(fun interface{}) error {
	funcType := reflect.TypeOf(fun)

	if funcType.Kind() != reflect.Func {
		err := errors.New(fmt.Sprintf("parse: should not register other than function"))
		if p.PanicOnError {
			panic(err.Error())
		}
		return err
	}

	if funcType.NumIn() != 1 || funcType.In(0).Kind() != reflect.String {
		err := errors.New(fmt.Sprintf("parse: field parser function should have one IN param with type string"))
		if p.PanicOnError {
			panic(err.Error())
		}
		return err
	}

	if funcType.NumOut() != 2 || funcType.Out(1).Name() != "error" {
		err := errors.New(fmt.Sprintf("parse: field parser function should have 2 out param with second type error"))
		if p.PanicOnError {
			panic(err.Error())
		}
		return err
	}

	fieldType := funcType.Out(0)

	p.dataTypeParser[fieldType] = reflect.ValueOf(fun)

	return nil
}

// parse the tsv file's first line, to get index of each field header
func (p *Parser) IndexByFieldHeader(titleLine string) map[string]int {
	fieldIndex := make(map[string]int, 0)
	for index, field := range strings.Split(titleLine, p.FieldsSplitter) {
		fieldIndex[field] = index
	}
	return fieldIndex
}

// parse a line of tsv file, to an user defined data struct, buffer provided by record
func (p *Parser) ParseAsOneRecord(line string, fieldIndex map[string]int, record interface{}) error {
	pointerType := reflect.TypeOf(record)

	if pointerType.Kind() != reflect.Ptr {
		err := errors.New("parse: record should be POINTOR to a table record struct")
		if p.PanicOnError {
			panic(err.Error())
		}
		return err
	}

	structType := pointerType.Elem()
	if structType.Kind() != reflect.Struct {
		err := errors.New("parse: record should be pointor to a table record STRUCT")
		if p.PanicOnError {
			panic(err.Error())
		}
		return err
	}

	fieldArrayData := strings.Split(line, p.FieldsSplitter)
	fieldArrayLen := len(fieldArrayData)

	structValue := reflect.ValueOf(record).Elem()

	//get and parse and set data of each field
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		header := field.Tag.Get(p.Tag)
		if header == "" {
			//we do not care about this filed
			continue
		}

		//get column index of this field, in order to get text data of this filed, as fieldArrayData[index]
		index, ok := fieldIndex[header]
		if !ok || index >= fieldArrayLen {
			err := errors.New(fmt.Sprintf("parse: header %s [%s] error %v", header, field.Name, fieldIndex))
			if p.PanicOnError {
				panic(err.Error())
				return err
			}
			continue
		}

		//get parse function of this data type for registry
		parseFunc, ok := p.dataTypeParser[field.Type]
		if !ok {
			err := errors.New(fmt.Sprintf("parse: can not find field parser for field %v", field.Type))
			if p.PanicOnError {
				panic(err.Error())
				return e
			}
			continue
		}

		//do parse the data to program data type
		out := parseFunc.Call([]reflect.Value{reflect.ValueOf(fieldArrayData[index])})
		value, err := out[0], out[1]
		if !err.IsNil() {
			ep := err.Interface().(error)
			e := errors.New(fmt.Sprintf("parse: parse field %s error: %s", field.Name, ep.Error()))
			if p.PanicOnError {
				panic(e.Error())
				return e
			}
			continue
		}

		//set value to the buffer field
		structValue.Field(i).Set(value)
	}

	return nil
}


