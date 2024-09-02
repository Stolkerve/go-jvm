package jvm

import (
	"bufio"
)

func ReadSection(javaClassFile *bufio.Reader, buffer []byte) error {
	if _, err := javaClassFile.Read(buffer); err != nil {
		return err
	}
	return nil
}
