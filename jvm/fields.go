package jvm

import (
	"bufio"
	"encoding/binary"
)

type FieldInfo struct {
	AccessFlags     AccessFlag
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

func ReadField(constantPool []*ConstantInfo, javaClassFile *bufio.Reader, sectionsReadBuffer []byte) (*FieldInfo, error) {
	if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
		return nil, err
	}
	var field FieldInfo
	field.AccessFlags = AccessFlag(binary.BigEndian.Uint16(sectionsReadBuffer))
	field.NameIndex = binary.BigEndian.Uint16(sectionsReadBuffer[2:])
	field.DescriptorIndex = binary.BigEndian.Uint16(sectionsReadBuffer[4:])
	field.AttributesCount = binary.BigEndian.Uint16(sectionsReadBuffer[6:])

	field.Attributes = make([]*AttributeInfo, field.AttributesCount)
	for i := range field.AttributesCount {
		attribute, err := ReadAttribute(constantPool, javaClassFile, make([]byte, 4))
		if err != nil {
			return nil, err
		}
		field.Attributes[i] = attribute
	}

	return &field, nil
}
