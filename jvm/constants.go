package jvm

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type ConstantPoolTag uint8

const (
	ConstantClassTag              ConstantPoolTag = 7
	ConstantFieldRefTag           ConstantPoolTag = 9
	ConstantMethodRefTag          ConstantPoolTag = 10
	ConstantInterfaceMethodRefTag ConstantPoolTag = 11
	ConstantStringTag             ConstantPoolTag = 8
	ConstantIntegerTag            ConstantPoolTag = 3
	ConstantFloatTag              ConstantPoolTag = 4
	ConstantLongTag               ConstantPoolTag = 5
	ConstantDoubleTag             ConstantPoolTag = 6
	ConstantNameAndTypeTag        ConstantPoolTag = 12
	ConstantUtf8Tag               ConstantPoolTag = 1
	ConstantMethodHandleTag       ConstantPoolTag = 15
	ConstantMethodTypeTag         ConstantPoolTag = 16
	ConstantInvokeDynamicTag      ConstantPoolTag = 18
)

func IsConstantPoolTag(v ConstantPoolTag) bool {
	switch v {
	case ConstantClassTag:
		return true
	case ConstantDoubleTag:
		return true
	case ConstantFieldRefTag:
		return true
	case ConstantFloatTag:
		return true
	case ConstantIntegerTag:
		return true
	case ConstantInterfaceMethodRefTag:
		return true
	case ConstantInvokeDynamicTag:
		return true
	case ConstantLongTag:
		return true
	case ConstantMethodHandleTag:
		return true
	case ConstantMethodTypeTag:
		return true
	case ConstantMethodRefTag:
		return true
	case ConstantNameAndTypeTag:
		return true
	case ConstantStringTag:
		return true
	case ConstantUtf8Tag:
		return true
	default:
		return false
	}
}

func (c ConstantPoolTag) String() string {
	switch c {
	case ConstantClassTag:
		return "ConstantClass"
	case ConstantDoubleTag:
		return "ConstantDouble"
	case ConstantFieldRefTag:
		return "ConstantFieldref"
	case ConstantFloatTag:
		return "ConstantFloat"
	case ConstantIntegerTag:
		return "ConstantInteger"
	case ConstantInterfaceMethodRefTag:
		return "ConstantInterfaceMethodref"
	case ConstantInvokeDynamicTag:
		return "ConstantInvokeDynamic"
	case ConstantLongTag:
		return "ConstantLong"
	case ConstantMethodHandleTag:
		return "ConstantMethodHandle"
	case ConstantMethodTypeTag:
		return "ConstantMethodType"
	case ConstantMethodRefTag:
		return "ConstantMethodref"
	case ConstantNameAndTypeTag:
		return "ConstantNameAndType"
	case ConstantStringTag:
		return "ConstantString"
	case ConstantUtf8Tag:
		return "ConstantUtf8"
	default:
		panic(fmt.Sprintf("unexpected main.ConstantPoolTags: %#v", c))
	}
}

type ConstantInfo struct {
	Tag  ConstantPoolTag
	Data interface{}
}

type ConstantClass struct {
	NameIndex uint16
}

type ConstantFieldRef struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

type ConstantMethodRef struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

type ConstantInterfaceMethodRef struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

type ConstantString struct {
	StringIndex uint16
}

type ConstantInteger = int32

type ConstantFloat = float32

type ConstantLong = int64

type ConstantDouble = float64

type ConstantNameAndType struct {
	NameIndex       uint16
	DescriptorIndex uint16
}

type ConstantUtf8 = string

type ConstantMethodHandle struct {
	ReferenceKind  uint8
	ReferenceIndex uint16
}

type ConstantMethodType struct {
	DescriptorIndex uint16
}

type ConstantInvokeDynamic struct {
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func ReadConstantPool(javaClassFile *bufio.Reader, sectionsReadBuffer []byte) (*ConstantInfo, error) {
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:1]); err != nil {
		return nil, err
	}

	tag := ConstantPoolTag(sectionsReadBuffer[:1][0])
	if !IsConstantPoolTag(tag) {
		return nil, fmt.Errorf("Invalid constant pool tag")
	}

	switch tag {
	case ConstantMethodRefTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: &ConstantMethodRef{
				ClassIndex:       binary.BigEndian.Uint16(sectionsReadBuffer),
				NameAndTypeIndex: binary.BigEndian.Uint16(sectionsReadBuffer[2:]),
			},
		}, nil
	case ConstantClassTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantClass{
				NameIndex: binary.BigEndian.Uint16(sectionsReadBuffer),
			},
		}, nil
	case ConstantDoubleTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		high := binary.BigEndian.Uint32(sectionsReadBuffer)

		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		low := binary.BigEndian.Uint32(sectionsReadBuffer)

		return &ConstantInfo{
			Tag:  tag,
			Data: ConstantDouble((int64(high) << 32) + int64(low)),
		}, nil
	case ConstantFieldRefTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantFieldRef{
				ClassIndex:       binary.BigEndian.Uint16(sectionsReadBuffer),
				NameAndTypeIndex: binary.BigEndian.Uint16(sectionsReadBuffer[2:]),
			},
		}, nil
	case ConstantFloatTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag:  tag,
			Data: ConstantFloat(binary.BigEndian.Uint32(sectionsReadBuffer)),
		}, nil
	case ConstantIntegerTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag:  tag,
			Data: ConstantInteger(binary.BigEndian.Uint32(sectionsReadBuffer)),
		}, nil
	case ConstantInterfaceMethodRefTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantInterfaceMethodRef{
				ClassIndex:       binary.BigEndian.Uint16(sectionsReadBuffer),
				NameAndTypeIndex: binary.BigEndian.Uint16(sectionsReadBuffer[2:]),
			},
		}, nil
	case ConstantInvokeDynamicTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantInvokeDynamic{
				BootstrapMethodAttrIndex: binary.BigEndian.Uint16(sectionsReadBuffer),
				NameAndTypeIndex:         binary.BigEndian.Uint16(sectionsReadBuffer[2:]),
			},
		}, nil
	case ConstantLongTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		high := binary.BigEndian.Uint32(sectionsReadBuffer)

		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		low := binary.BigEndian.Uint32(sectionsReadBuffer)

		return &ConstantInfo{
			Tag:  tag,
			Data: ConstantLong((int64(high) << 32) + int64(low)),
		}, nil
	case ConstantMethodHandleTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:3]); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantMethodHandle{
				ReferenceKind:  sectionsReadBuffer[0],
				ReferenceIndex: binary.BigEndian.Uint16(sectionsReadBuffer[1:]),
			},
		}, nil
	case ConstantMethodTypeTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantMethodType{
				DescriptorIndex: binary.BigEndian.Uint16(sectionsReadBuffer),
			},
		}, nil
	case ConstantNameAndTypeTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantNameAndType{
				NameIndex:       binary.BigEndian.Uint16(sectionsReadBuffer),
				DescriptorIndex: binary.BigEndian.Uint16(sectionsReadBuffer[2:]),
			},
		}, nil
	case ConstantStringTag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
			return nil, err
		}
		return &ConstantInfo{
			Tag: tag,
			Data: ConstantString{
				StringIndex: binary.BigEndian.Uint16(sectionsReadBuffer),
			},
		}, nil
	case ConstantUtf8Tag:
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
			return nil, err
		}
		lenght := binary.BigEndian.Uint16(sectionsReadBuffer)

		ut8Buffer := make([]byte, lenght)
		if err := ReadSection(javaClassFile, ut8Buffer); err != nil {
			return nil, err
		}

		return &ConstantInfo{
			Tag:  tag,
			Data: ConstantUtf8(string(ut8Buffer)),
		}, nil
	default:
		panic(fmt.Sprintf("unexpected main.ConstantPoolTag: %#v", tag))
	}
}
