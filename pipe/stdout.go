package pipe

import "fmt"

func WriteToStdOut(bytes []byte) (int, error) {

	fmt.Println(bytes)
	return len(bytes), nil

}
