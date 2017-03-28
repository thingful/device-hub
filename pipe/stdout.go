package pipe

import "fmt"

func WriteToStdOut(bytes []byte) (int, error) {

	fmt.Println(string(bytes))
	return len(bytes), nil

}
