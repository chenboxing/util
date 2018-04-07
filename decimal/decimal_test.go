package decimal

import "testing"
import "fmt"

func TestAdd(t *testing.T) {
	res := Add(1, 2)
	fmt.Println(res.Decimal())
}

func TestSub(t *testing.T) {
	res := Sub(1, 2, 3)
	fmt.Println(res.Decimal())
}
func TestMul(t *testing.T) {
	res := Mul(1, 2, 3)
	fmt.Println(res.Decimal())
}
func TestDiv(t *testing.T) {
	res := Div(1, 2, 3)
	fmt.Println(res.Decimal())
}

func TestDecimal(t *testing.T) {
	res := Add(0, 1)
	for i := 1; i < 100000000000; i++ {
		res = Mul(res, i)
		fmt.Printf("%v\n", res.Decimal())
	}

	fmt.Println(res)
}
