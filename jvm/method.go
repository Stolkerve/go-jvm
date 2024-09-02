package jvm

import (
	"bufio"
	"encoding/binary"
)

type MethodInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	Name            string
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

func ReadMethod(constantPool []*ConstantInfo, javaClassFile *bufio.Reader, sectionsReadBuffer []byte) (*MethodInfo, error) {
	if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
		return nil, err
	}
	var method MethodInfo
	method.AccessFlags = binary.BigEndian.Uint16(sectionsReadBuffer)
	method.NameIndex = binary.BigEndian.Uint16(sectionsReadBuffer[2:])
	method.DescriptorIndex = binary.BigEndian.Uint16(sectionsReadBuffer[4:])
	method.AttributesCount = binary.BigEndian.Uint16(sectionsReadBuffer[6:])

	method.Name = constantPool[method.NameIndex-1].Data.(ConstantUtf8)

	method.Attributes = make([]*AttributeInfo, method.AttributesCount)
	for i := range method.AttributesCount {
		attribute, err := ReadAttribute(constantPool, javaClassFile, make([]byte, 4))
		if err != nil {
			return nil, err
		}
		method.Attributes[i] = attribute
	}

	return &method, nil
}
