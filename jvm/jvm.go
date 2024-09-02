package jvm

import (
	"bufio"
	"fmt"
	"os"
)

type StackType int

const (
	StackTypeInt = StackType(iota)
	StackTypeConstant
	StackTypeStaticClass
)

type StackData struct {
	Type StackType
	Data interface{}
}

type Jvm struct {
	Class *JavaClass
}

func NewJvm(filename string) (*Jvm, error) {
	javaClassFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer javaClassFile.Close()

	class, err := NewJavaClass(bufio.NewReader(javaClassFile))
	if err != nil {
		return nil, err
	}

	return &Jvm{
		Class: class,
	}, nil
}

func RunJvm(jvm *Jvm) {
	// fmt.Println(jvm.Class)
	var mainMethod *MethodInfo

	for _, m := range jvm.Class.Methods {
		if m.Name == "main" {
			mainMethod = m
		}
	}

	mainCodeAttr := mainMethod.Attributes[0].Data.(CodeAttribute)

	fmt.Println("Running", mainMethod.Name, "function code")
	RunCodeAttr(jvm, &mainCodeAttr)
}

func RunCodeAttr(jvm *Jvm, codeAttr *CodeAttribute) {
	pc := 0
	code := codeAttr.Code

	// fmt.Printf("%X\n", code)
	locals := make([]StackData, codeAttr.MaxLocals)
	stack := make([]StackData, 0, codeAttr.MaxStack)
	for pc < len(code) {
		switch code[pc] {
		case 0x10: // Push byte
			pc += 1
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(code[pc]),
			})
		case 0x3C: // istore_1 Store int into local variable
			v := stack[len(stack)-1]
			if v.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "At instruction 0x%X expected stack value of type int\n", code[pc])
				break
			}
			locals[1] = v
			stack = stack[:len(stack)-1]
		case 0x3D: // istore_2 Store int into local variable
			v := stack[len(stack)-1]
			if v.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "At instruction 0x%X expected stack value of type int\n", code[pc])
				break
			}
			locals[2] = v
			stack = stack[:len(stack)-1]
		case 0x03: // iconst_0
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(0),
			})
		case 0x04: // iconst_1
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(1),
			})
		case 0x05: // iconst_2
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(2),
			})
		case 0x06: // iconst_3
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(3),
			})
		case 0x60: // iadd
			v1 := stack[len(stack)-1]
			v2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if v1.Type != StackTypeInt && v2.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "At instruction 0x%X expected stack value of type int\n", code[pc])
			}

			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(v1.Data.(int32) + v2.Data.(int32)),
			})
		case 0x68: // imul
			v1 := stack[len(stack)-1]
			v2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if v1.Type != StackTypeInt && v2.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "At instruction 0x%X expected stack value of type int\n", code[pc])
			}

			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int32(v1.Data.(int32) * v2.Data.(int32)),
			})
		case 0xB2: // Get static field from class
			pc += 1
			indexByte1 := uint16(code[pc])
			pc += 1
			indexByte2 := uint16(code[pc])
			index := (indexByte1 << 8) | indexByte2
			fieldRef := jvm.Class.ConstantPool[index-1].Data.(ConstantFieldRef)
			class := jvm.Class.ConstantPool[fieldRef.ClassIndex-1].Data.(ConstantClass)
			className := jvm.Class.ConstantPool[class.NameIndex-1].Data.(ConstantUtf8)
			method := jvm.Class.ConstantPool[fieldRef.NameAndTypeIndex-1].Data.(ConstantNameAndType)
			methodName := jvm.Class.ConstantPool[method.NameIndex-1].Data.(ConstantUtf8)
			if className == "java/lang/System" && methodName == "out" {
				stack = append(stack, StackData{
					Type: StackTypeStaticClass,
					Data: "JavaPrintStream",
				})
			} else {
				fmt.Fprintf(os.Stderr, "Unsupported static class %s method %s\n", className, methodName)
				break
			}
		case 0x12: // Push item from run-time constant pool
			pc += 1
			index := uint16(code[pc])
			stack = append(stack, StackData{
				Type: StackTypeConstant,
				Data: index,
			})
		case 0x1B: // iload_1 load an int value from local variable 1
			int := locals[1]
			if int.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "Expected value type int\n")
				break
			}
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int.Data,
			})
		case 0x1c: // iload_2 load an int value from local variable 2
			int := locals[2]
			if int.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "Expected value type int\n")
				break
			}
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int.Data,
			})
		case 0x1d: // iload_3 load an int value from local variable 2
			int := locals[3]
			if int.Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "Expected value type int\n")
				break
			}
			stack = append(stack, StackData{
				Type: StackTypeInt,
				Data: int.Data,
			})
		case 0x84: //iinc  	increment local variable #index by signed byte const
			pc += 1
			index := int(code[pc])
			pc += 1
			v := int32(code[pc])

			if locals[index].Type != StackTypeInt {
				fmt.Fprintf(os.Stderr, "Expected value type int\n")
				break
			}
			locals[index].Data = locals[index].Data.(int32) + v
		case 0xB6: // Invoke instance method; dispatch based on class
			pc += 1
			indexByte1 := uint16(code[pc])
			pc += 1
			indexByte2 := uint16(code[pc])
			index := (indexByte1 << 8) | indexByte2

			methodRef := jvm.Class.ConstantPool[index-1].Data.(*ConstantMethodRef)

			class := jvm.Class.ConstantPool[methodRef.ClassIndex-1].Data.(ConstantClass)
			className := jvm.Class.ConstantPool[class.NameIndex-1].Data.(ConstantUtf8)
			method := jvm.Class.ConstantPool[methodRef.NameAndTypeIndex-1].Data.(ConstantNameAndType)
			methodName := jvm.Class.ConstantPool[method.NameIndex-1].Data.(ConstantUtf8)

			if className == "java/io/PrintStream" && methodName == "println" {
				if len(stack) < 2 {
					fmt.Fprintf(os.Stderr, "expected two arguments in class %s on method %s, found %d\n", className, methodName, len(stack))
					break
				}
				value := stack[len(stack)-1]
				class := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				if class.Type == StackTypeStaticClass {
					if class.Data.(string) != "JavaPrintStream" {
						fmt.Fprintf(os.Stderr, "expected %s class, found %s\n", "JavaPrintStream", class.Data.(string))
						break
					}
				}
				switch value.Type {
				case StackTypeConstant:
					constant := jvm.Class.ConstantPool[value.Data.(uint16)-1]
					switch constant.Tag {
					case ConstantStringTag:
						str := jvm.Class.ConstantPool[constant.Data.(ConstantString).StringIndex-1].Data.(ConstantUtf8)
						fmt.Println(str)
					}
				default:
				case StackTypeInt:
					fmt.Println(value.Data.(int32))
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unsupported class %s method %s\n", className, methodName)
				break
			}

		case 0xB1:
			break
		default:
			fmt.Fprintf(os.Stderr, "opcode 0x%02X not supported\n", code[pc])
		}
		pc += 1
		// fmt.Println(stack)
	}
}
