package  main
 
import "fmt"

func main(){
	fmt.Println("go" + "lang")
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)
	fmt.Println("0.1+0.2 =", 0.1+0.2)
	fmt.Println(true && false)
	fmt.Println(true || false)


	f1 := float64(0.1)
	f2 := float64(0.2)
	f3 := f1 + f2
	fmt.Println(f3) // 0.30000000000000004
	
}