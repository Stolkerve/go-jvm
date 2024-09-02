package jvm

import (
	"fmt"
	"strings"
)

func ParseBaseType(c byte) string {
	switch c {
	case 'B':
		return "byte"
	case 'C':
		return "char"
	case 'D':
		return "double"
	case 'F':
		return "float"
	case 'I':
		return "int"
	case 'J':
		return "long"
	case 'Z':
		return "boolean"
	case 'V':
		return "void"
	}
	return ""
}

func ParseFieldType(descriptor string, pos *int) string {
	bt := ParseBaseType(descriptor[*pos])
	if len(bt) != 0 {
		return bt
	}
	if descriptor[*pos] == '[' {
		*pos += 1
		return fmt.Sprintf("%s[]", ParseFieldType(descriptor, pos))
	} else if descriptor[*pos] == 'L' {
		semiColPos := strings.IndexByte(descriptor[*pos:], ';')
		return strings.ReplaceAll(descriptor[*pos+1:semiColPos], "/", ".")
	}
	return ""
}

func ParseDescriptor(descriptor, name string) string {
	isMethod := strings.HasPrefix(descriptor, "(")
	if isMethod {
		returnType := descriptor[len(descriptor)-1]
		params := descriptor[1 : len(descriptor)-2]
		paramsTypes := make([]string, 0)
		if len(params) != 0 {
			pos := 0
			for pos < len(params) {
				t := ParseFieldType(params, &pos)
				pos += 1

				if len(t) == 0 {
					break
				}
				paramsTypes = append(paramsTypes, t)
			}
		}
		return fmt.Sprintf("%s %s(%s)", ParseBaseType(returnType), name, strings.Join(paramsTypes, ", "))
	}
	pos := 0
	return fmt.Sprintf("%s %s", ParseFieldType(descriptor, &pos), name)
}
