package jvm

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

type AccessFlag uint16

const (
	AccPublic     AccessFlag = 0x0001
	AccPrivate    AccessFlag = 0x0002
	AccProtected  AccessFlag = 0x0004
	AccStatic     AccessFlag = 0x0008
	AccFinal      AccessFlag = 0x0010
	AccSuper      AccessFlag = 0x0020
	AccVolatile   AccessFlag = 0x0040
	AccTransient  AccessFlag = 0x0080
	AccNative     AccessFlag = 0x0100
	AccInterface  AccessFlag = 0x0200
	AccAbstract   AccessFlag = 0x0400
	AccStrict     AccessFlag = 0x0800
	AccSynthetic  AccessFlag = 0x1000
	AccAnnotation AccessFlag = 0x2000
	AccEnum       AccessFlag = 0x4000
)

func (f AccessFlag) String() string {
	switch f {
	case AccAbstract:
		return "abstract"
	case AccAnnotation:
		return "annotation"
	case AccEnum:
		return "enum"
	case AccFinal:
		return "final"
	case AccInterface:
		return "interface"
	case AccNative:
		return "native"
	case AccPrivate:
		return "private"
	case AccProtected:
		return "protected"
	case AccPublic:
		return "public"
	case AccStatic:
		return "static"
	case AccStrict:
		return "strict"
	case AccSuper:
		return "super"
	case AccSynthetic:
		return "synthetic"
	case AccTransient:
		return "transient"
	case AccVolatile:
		return "volatile"
	default:
		panic(fmt.Sprintf("unexpected jvm.AccessFlag: %#v", f))
	}
}

type JavaClass struct {
	Version      string
	ConstantPool []*ConstantInfo
	Interfaces   []uint16
	Fields       []*FieldInfo
	Methods      []*MethodInfo
	Attributes   []*AttributeInfo
	AccessFlags  AccessFlag
	ThisClass    uint16
	SuperClass   uint16
}

func NewJavaClass(javaClassFile *bufio.Reader) (*JavaClass, error) {
	var javaClass JavaClass

	// Decode magic number
	var sectionsReadBuffer = make([]byte, 4)
	if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
		return nil, err
	}
	magicNumber := binary.BigEndian.Uint32(sectionsReadBuffer)
	if magicNumber != 0xCAFEBABE {
		return nil, fmt.Errorf("Give me CAFEBABE")
	}

	// Decode minor and major version
	if err := ReadSection(javaClassFile, sectionsReadBuffer); err != nil {
		return nil, err
	}
	minorVersion := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	majorVersion := binary.BigEndian.Uint16(sectionsReadBuffer[2:])
	javaClass.Version = fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	// Decode constant pool count
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	constantPoolCount := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	javaClass.ConstantPool = make([]*ConstantInfo, constantPoolCount-1)

	// Decode constant pool
	for i := range constantPoolCount - 1 {
		constant, err := ReadConstantPool(javaClassFile, make([]byte, 4))
		if err != nil {
			return nil, err
		}
		javaClass.ConstantPool[i] = constant
		// fmt.Println(i+1, constant.Tag.String())
	}

	// Decode access flags
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	javaClass.AccessFlags = AccessFlag(binary.BigEndian.Uint16(sectionsReadBuffer[:2]))

	// Decode this class
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	javaClass.ThisClass = binary.BigEndian.Uint16(sectionsReadBuffer[:2])

	// Decode super class
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	javaClass.SuperClass = binary.BigEndian.Uint16(sectionsReadBuffer[:2])

	// Decode interfaces count
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	interfacesCount := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	javaClass.Interfaces = make([]uint16, interfacesCount)
	for i := range interfacesCount {
		if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
			return nil, err
		}
		javaClass.Interfaces[i] = binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	}

	// Decode fields count
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	fieldsCount := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	javaClass.Fields = make([]*FieldInfo, fieldsCount)
	for i := range fieldsCount {
		field, err := ReadField(javaClass.ConstantPool, javaClassFile, make([]byte, 8))
		if err != nil {
			return nil, err
		}
		javaClass.Fields[i] = field
	}

	// Decode methods count
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	methodsCount := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	javaClass.Methods = make([]*MethodInfo, methodsCount)
	for i := range methodsCount {
		method, err := ReadMethod(javaClass.ConstantPool, javaClassFile, make([]byte, 8))
		if err != nil {
			return nil, err
		}
		javaClass.Methods[i] = method
	}

	// Decode attributes count
	if err := ReadSection(javaClassFile, sectionsReadBuffer[:2]); err != nil {
		return nil, err
	}
	attributesCount := binary.BigEndian.Uint16(sectionsReadBuffer[:2])
	javaClass.Attributes = make([]*AttributeInfo, attributesCount)
	for i := range attributesCount {
		attribute, err := ReadAttribute(javaClass.ConstantPool, javaClassFile, make([]byte, 4))
		if err != nil {
			return nil, err
		}
		javaClass.Attributes[i] = attribute
	}

	return &javaClass, nil
}

func (c *JavaClass) String() string {
	var output bytes.Buffer
	fmt.Fprintf(&output, "Version: %s\n", c.Version)
	fmt.Fprintf(&output, "Access flag: (0x%X) %s\n", c.AccessFlags, c.AccessFlags)

	thisClass := c.ConstantPool[c.ThisClass-1].Data.(ConstantClass)
	thisClassName := c.ConstantPool[thisClass.NameIndex-1].Data.(ConstantUtf8)
	fmt.Fprintf(&output, "This class: (#%d) %s\n", c.ThisClass, thisClassName)

	superClass := c.ConstantPool[c.SuperClass-1].Data.(ConstantClass)
	superClassName := c.ConstantPool[superClass.NameIndex-1].Data.(ConstantUtf8)
	fmt.Fprintf(&output, "Super class: (#%d) %s\n", c.SuperClass, superClassName)

	fmt.Fprintf(&output, "Interfaces: (%d)\n", len(c.Interfaces))
	fmt.Fprintf(&output, "Fields: (%d)\n", len(c.Fields))
	fmt.Fprintf(&output, "Methods: (%d)\n", len(c.Methods))

	for i, m := range c.Methods {
		fmt.Fprintf(&output, "\t#%d %s: \n", i+1, m.Name)
	}

	fmt.Fprintf(&output, "Constant pool: (%d)\n", len(c.ConstantPool))
	for i, constant := range c.ConstantPool {
		fmt.Fprintf(&output, "\t#%d %s: ", i+1, constant.Tag)
		switch constant.Tag {
		case ConstantClassTag:
			class := constant.Data.(ConstantClass)
			fmt.Fprintf(&output, "#%d %s", class.NameIndex, c.ConstantPool[class.NameIndex-1].Data.(ConstantUtf8))
		case ConstantDoubleTag:
			double := constant.Data.(ConstantDouble)
			fmt.Fprintf(&output, "%f", double)
		case ConstantFieldRefTag:
			fieldRef := constant.Data.(ConstantFieldRef)
			fmt.Fprintf(&output, "#%d #%d", fieldRef.ClassIndex, fieldRef.NameAndTypeIndex)
		case ConstantFloatTag:
			float := constant.Data.(ConstantFloat)
			fmt.Fprintf(&output, "%f", float)
		case ConstantIntegerTag:
			int := constant.Data.(ConstantInteger)
			fmt.Fprintf(&output, "#%d", int)
		case ConstantInterfaceMethodRefTag:
			interfaceMethodRefTag := constant.Data.(ConstantInterfaceMethodRef)
			fmt.Fprintf(&output, "#%d #%d", interfaceMethodRefTag.ClassIndex, interfaceMethodRefTag.NameAndTypeIndex)
		case ConstantInvokeDynamicTag:
			invokeDynamic := constant.Data.(ConstantInvokeDynamic)
			fmt.Fprintf(&output, "#%d #%d", invokeDynamic.BootstrapMethodAttrIndex, invokeDynamic.NameAndTypeIndex)
		case ConstantLongTag:
			long := constant.Data.(ConstantLong)
			fmt.Fprintf(&output, "%d", long)
		case ConstantMethodHandleTag:
			methodHandle := constant.Data.(ConstantMethodHandle)
			fmt.Fprintf(&output, "#%d #%d", methodHandle.ReferenceKind, methodHandle.ReferenceIndex)
		case ConstantMethodRefTag:
			methodRef := constant.Data.(*ConstantMethodRef)
			className := c.ConstantPool[c.ConstantPool[methodRef.ClassIndex-1].Data.(ConstantClass).NameIndex-1].Data.(ConstantUtf8)
			nameAndType := c.ConstantPool[methodRef.NameAndTypeIndex-1].Data.(ConstantNameAndType)
			name := c.ConstantPool[nameAndType.NameIndex-1].Data.(ConstantUtf8)
			descriptor := c.ConstantPool[nameAndType.DescriptorIndex-1].Data.(ConstantUtf8)
			fmt.Fprintf(&output, "#%d #%d %s %s", methodRef.ClassIndex, methodRef.NameAndTypeIndex, className, ParseDescriptor(descriptor, name))
		case ConstantMethodTypeTag:
			methodType := constant.Data.(ConstantMethodType)
			fmt.Fprintf(&output, "#%d", methodType.DescriptorIndex)
		case ConstantNameAndTypeTag:
			nameAndType := constant.Data.(ConstantNameAndType)
			name := c.ConstantPool[nameAndType.NameIndex-1].Data.(ConstantUtf8)
			descriptor := c.ConstantPool[nameAndType.DescriptorIndex-1].Data.(ConstantUtf8)
			fmt.Fprintf(&output, "#%d #%d %s", nameAndType.NameIndex, nameAndType.DescriptorIndex, ParseDescriptor(descriptor, name))
		case ConstantStringTag:
			string := constant.Data.(ConstantString)
			fmt.Fprintf(&output, "#%d %s", string.StringIndex, c.ConstantPool[string.StringIndex-1].Data.(ConstantUtf8))
		case ConstantUtf8Tag:
			fmt.Fprintf(&output, "%s", constant.Data.(ConstantUtf8))
		}
		fmt.Fprintf(&output, "\n")
	}

	fmt.Fprintf(&output, "Attributes: (%d)\n", len(c.Attributes))
	for i, attr := range c.Attributes {
		fmt.Fprintf(&output, "\t#%d %s: ", i+1, attr.AttributeType)
		switch attr.AttributeType {
		case SourceFileAttr:
			source := attr.Data.(SourceFileAttribute)
			fmt.Fprintf(&output, "%s", source.Sourcefile)
		}
		fmt.Fprintf(&output, "\n")
	}

	return output.String()
}
