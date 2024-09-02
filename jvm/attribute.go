package jvm

import (
	"bufio"
	"bytes"

	"encoding/binary"
	"fmt"
)

type AttributeType string

const (
	ConstantValueAttr                        AttributeType = "ConstantValue"
	CodeAttr                                 AttributeType = "Code"
	StackMapTableAttr                        AttributeType = "StackMapTable"
	ExceptionsAttr                           AttributeType = "Exceptions"
	InnerClassesAttr                         AttributeType = "InnerClasses"
	EnclosingMethodAttr                      AttributeType = "EnclosingMethod"
	SyntheticAttr                            AttributeType = "Synthetic"
	SignatureAttr                            AttributeType = "Signature"
	SourceFileAttr                           AttributeType = "SourceFile"
	SourceDebugExtensionAttr                 AttributeType = "SourceDebugExtension"
	LineNumberTableAttr                      AttributeType = "LineNumberTable"
	LocalVariableTableAttr                   AttributeType = "LocalVariableTable"
	LocalVariableTypeTableAttr               AttributeType = "LocalVariableTypeTable"
	DeprecatedAttr                           AttributeType = "Deprecated"
	RuntimeVisibleAnnotationsAttr            AttributeType = "RuntimeVisibleAnnotations"
	RuntimeInvisibleAnnotationsAttr          AttributeType = "RuntimeInvisibleAnnotations"
	RuntimeVisibleParameterAnnotationsAttr   AttributeType = "RuntimeVisibleParameterAnnotations"
	RuntimeInvisibleParameterAnnotationsAttr AttributeType = "RuntimeInvisibleParameterAnnotations"
	AnnotationDefaultAttr                    AttributeType = "AnnotationDefault"
	BootstrapMethodsAttr                     AttributeType = "BootstrapMethods"
)

func IsAttributeType(attr AttributeType) bool {
	switch attr {
	case AnnotationDefaultAttr:
		return true
	case BootstrapMethodsAttr:
		return true
	case CodeAttr:
		return true
	case ConstantValueAttr:
		return true
	case DeprecatedAttr:
		return true
	case EnclosingMethodAttr:
		return true
	case ExceptionsAttr:
		return true
	case InnerClassesAttr:
		return true
	case LineNumberTableAttr:
		return true
	case LocalVariableTableAttr:
		return true
	case LocalVariableTypeTableAttr:
		return true
	case RuntimeInvisibleAnnotationsAttr:
		return true
	case RuntimeInvisibleParameterAnnotationsAttr:
		return true
	case RuntimeVisibleAnnotationsAttr:
		return true
	case RuntimeVisibleParameterAnnotationsAttr:
		return true
	case SignatureAttr:
		return true
	case SourceDebugExtensionAttr:
		return true
	case SourceFileAttr:
		return true
	case StackMapTableAttr:
		return true
	case SyntheticAttr:
		return true
	default:
		return false
	}
}

type CodeAttribute struct {
	MaxStack        uint16
	MaxLocals       uint16
	CodeLength      uint32
	Code            []byte
	ExceptionsTable []struct {
		StartPc   uint16
		EndPc     uint16
		HandlerPc uint16
		CatchType uint16
	}
	Attributes []*AttributeInfo
}

type SourceFileAttribute struct {
	SourcefileIndex uint16
	Sourcefile      string
}

type AttributeInfo struct {
	AttributeNameIndex uint16
	AttributeType      AttributeType
	Data               interface{}
}

type LineNumberTableAttributeData struct {
	StartPc    uint16
	LineNumber uint16
}

type LineNumberTableAttribute []LineNumberTableAttributeData

func ReadAttribute(constantPool []*ConstantInfo, javaClassFile *bufio.Reader, readBuffer []byte) (*AttributeInfo, error) {
	if err := ReadSection(javaClassFile, readBuffer[:2]); err != nil {
		return nil, err
	}
	var attribute AttributeInfo
	attribute.AttributeNameIndex = binary.BigEndian.Uint16(readBuffer)

	attribute.AttributeType = AttributeType(constantPool[attribute.AttributeNameIndex-1].Data.(ConstantUtf8))
	if !IsAttributeType(attribute.AttributeType) {
		return nil, fmt.Errorf("unknow attribute type %s", attribute.AttributeType)
	}

	if err := ReadSection(javaClassFile, readBuffer); err != nil {
		return nil, err
	}
	attributeLength := binary.BigEndian.Uint32(readBuffer)
	info := make([]byte, attributeLength)
	if _, err := javaClassFile.Read(info); err != nil {
		return nil, err
	}

	switch attribute.AttributeType {
	case AnnotationDefaultAttr:
	case BootstrapMethodsAttr:
	case CodeAttr:
		codeAttr := CodeAttribute{}
		codeAttr.MaxStack = binary.BigEndian.Uint16(info)
		codeAttr.MaxLocals = binary.BigEndian.Uint16(info[2:])
		codeAttr.CodeLength = binary.BigEndian.Uint32(info[4:])
		codeAttr.Code = info[8 : 8+codeAttr.CodeLength]

		offset := 8 + codeAttr.CodeLength
		exceptionTableLength := binary.BigEndian.Uint16(info[offset:])
		offset += 2 + (uint32(exceptionTableLength) * 8)
		attributesCount := binary.BigEndian.Uint16(info[offset:])
		offset += 2

		codeAttr.Attributes = make([]*AttributeInfo, attributesCount)
		for i := range attributesCount {
			attrBuff := bufio.NewReader(bytes.NewBuffer(info[offset:]))
			attr, err := ReadAttribute(constantPool, attrBuff, make([]byte, 4))
			if err != nil {
				return nil, err
			}
			codeAttr.Attributes[i] = attr
		}

		attribute.Data = codeAttr
	case ConstantValueAttr:
	case DeprecatedAttr:
	case EnclosingMethodAttr:
	case ExceptionsAttr:
	case InnerClassesAttr:
	case LineNumberTableAttr:
		lenght := binary.BigEndian.Uint16(info)
		lineNumberTable := make(LineNumberTableAttribute, lenght)
		for i := range lenght {
			lineNumberTable[int(i)] = LineNumberTableAttributeData{
				StartPc:    binary.BigEndian.Uint16(info[(i*4)+2:]),
				LineNumber: binary.BigEndian.Uint16(info[(i*4)+2+2:]),
			}
		}

		attribute.Data = lineNumberTable
	case LocalVariableTableAttr:
	case LocalVariableTypeTableAttr:
	case RuntimeInvisibleAnnotationsAttr:
	case RuntimeInvisibleParameterAnnotationsAttr:
	case RuntimeVisibleAnnotationsAttr:
	case RuntimeVisibleParameterAnnotationsAttr:
	case SignatureAttr:
	case SourceDebugExtensionAttr:
	case SourceFileAttr:
		sourceAttr := SourceFileAttribute{}
		sourceAttr.SourcefileIndex = binary.BigEndian.Uint16(info)
		sourceAttr.Sourcefile = constantPool[sourceAttr.SourcefileIndex-1].Data.(ConstantUtf8)
		attribute.Data = sourceAttr
	case StackMapTableAttr:
	case SyntheticAttr:
	}

	return &attribute, nil
}
