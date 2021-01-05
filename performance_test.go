package lua

import (
	"fmt"
	"testing"
	"time"
)

func TestPerformanceLua(t *testing.T) {
	var luaCode = "local function fib(n)\n if n < 2 then return n end\n    return fib(n - 2) + fib(n - 1)\nend\n\n for i = 1, 10 do fib(35) \nend"
	fmt.Println("Start golua-fib")
	start := time.Now()
	L := NewState()
	defer L.Close()
	if err := L.DoString(luaCode); err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(duration.Milliseconds())
}
