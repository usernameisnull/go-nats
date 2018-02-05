//go build -ldflags "-H windowsgui" no-window.go
package main
import (
	"fmt"
	"time"
)

func main(){
	fmt.Println("Hello,world!")
	time.Sleep(time.Second*10)
}
