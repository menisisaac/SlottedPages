package Structure

import (
	"fmt"
	"testing"
)

func TestFirst(t *testing.T) {
	l := concatenate(5, 6)
	fmt.Println(int(second(l)))
}

func TestFileManager(t *testing.T) {
	m := CreateFileManager(2048)
	l := make([]int64, 35)
	for x := 0; x < 35; x += 1 {
		l[x] = (m.add(0, x))
	}
	for x := 0; x < 35; x += 1 {
		fmt.Println(m.get(0, l[x]))
	}
}
