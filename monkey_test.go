package monkey

import (
	"fmt"
	"testing"
)

func TestPrintln(t *testing.T) {
	fmt.Println("hello")

	patch := func(a ...interface{}) (int, error) {
		return fmt.Print("what!!\n")
	}
	Patch(fmt.Println, patch)

	fmt.Println("hello") // print "what!!"
	fmt.Println("?????") // print "what!!"
}
