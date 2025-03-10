package btils

import (
	"math/rand"
	"time"
)

func main() {
	tm := NewThreadManager[string](2, func(in string) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		println("Handled", in)
	})

	tm.Start()

	for _, str := range []string{"Foo", "Baar", "Baloo", "Golang"} {
		tm.Feed(str)
	}

	for !tm.IsDone() {
		time.Sleep(10 * time.Millisecond)
		println("Waiting ...")
	}

	tm.Stop()
	println("Done.")
}
