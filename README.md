# btils

A collection of commonly used utility and quality-of-life functions for Go projects. **btils** includes a lightweight concurrency worker pool, a fast (non-cryptographic) unique identifier generator, helper functions for conditional logic and default values, and JSON utility functions optimized with [goccy/go-json](https://github.com/goccy/go-json).

---

## Contents

- [Threader](#threader)
- [UID Utilities](#uid-utilities)
- [Fast Random Number Generation](#fast-random-number-generation)
- [Quality-of-Life Helpers](#quality-of-life-helpers)
- [JSON Utilities](#json-utilities)

---

## Threader

The **Threader** provides a simplistic worker pool structure to efficiently handle concurrent tasks. It is generic over any type `T` and allows you to feed tasks to a pool of workers that process items concurrently.

### How It Works

- **Creation:**  
  Create a new thread manager using `NewThreadManager[T](workers int, callback func(in T))`. The `workers` parameter determines the number of concurrent goroutines and `callback` is the function that processes each task.

- **Feeding Tasks:**  
  Use `Feed(in T)` to send tasks to the worker pool. Internally, an atomic counter tracks the number of tasks.

- **Monitoring:**  
  `IsDone()` checks if all tasks have been processed (i.e. the counter is 0).

- **Stopping:**  
  When done, call `Stop()` to close the underlying channel and terminate the worker goroutines.

### When to use

The **Threader** is ideal to use when the individual tasks take a non-predictable amount of time to complete. Due to the **Threader**s architecture, it will distribute the work as fast as possible across all workers. Whereas similar design patterns may result in threads idling while there is still work to do

### Example Usage

```go
package main

import (
	"math/rand"
	"time"

	"github.com/41Baloo/btils"
)

func main() {
	// Create a new thread manager with 2 workers that process strings.
	tm := btils.NewThreadManager[string](2, func(in string) {
		// Random delay so the next free worker is random
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		println("Handled", in)
	})

	// Start the worker pool.
	tm.Start()

	// Feed tasks to the pool.
	for _, str := range []string{"Foo", "Baar", "Baloo", "Golang"} {
		tm.Feed(str)
	}

	// Wait until all tasks are complete.
	for !tm.IsDone() {
		time.Sleep(10 * time.Millisecond)
		println("Waiting ...")
	}

	// Stop the worker pool.
	tm.Stop()
	println("Done.")
}
```

#### Example Output

```
Waiting ...
Waiting ...
Waiting ...
Handled Baar
Waiting ...
Waiting ...
Handled Foo
Waiting ...
Waiting ...
Waiting ...
Handled Baloo
Waiting ...
Waiting ...
Handled Golang
Done.
```

---

## UID Utilities

The UID utilities provide a simple 16-byte unique identifier. **Note:** This UID is **not** RFC4122 compliant and should **not** be used for cryptographic purposes.

### Features

- **UID Type:**  
  A UID is defined as a `[16]byte` array.

- **Conversion:**  
  - `UIDFromString(s string) *UID` converts a string (of at least 16 characters) into a UID.
  - `ToString()` returns the UID as a string.

- **Validation:**  
  `IsValid()` checks if the UID contains only allowed characters (alphanumeric, underscore, and dash).

- **Generation:**  
  `NewUID(b *UID)` rapidly generates a new UID using the fast random number generator.  
  *It reuses old UIDs if desired and uses low-level unsafe conversions for speed.*

### Example

```go
package main

import (
	"fmt"

	"github.com/41Baloo/btils"
)

func main() {
	var uid btils.UID
	btils.NewUID(&uid)

	uidStr := uid.ToString()
	fmt.Println("New UID:", uidStr)

	// Validate UID
	if uid.IsValid() {
		fmt.Println("The UID is valid.")
	} else {
		fmt.Println("The UID is NOT valid.")
	}
}
```

---

## Fast Random Number Generation

The package provides a fast random number function:

- **Fastrand:**  
  The function `Fastrand()` is linked to Go's internal `runtime.cheaprand` and provides fast (but not cryptographically secure) random numbers. It is used internally by `NewUID`.

---

## Quality-of-Life Helpers

### None

`None[T any]() T` returns the zero value for the type `T`.  
Useful when you need to initialize a generic variable without knowing its type.

### If

`If[T any](cond bool, truely, falsely T) T` acts as a ternary operator. It returns `truely` if `cond` is `true`, and `falsely` otherwise.

### Example

```go
package main

import (
	"fmt"

	"github.com/41Baloo/btils"
)

func main() {
	defaultInt := btils.None[int]()
	fmt.Println("Default int value:", defaultInt) // 0

	result := btils.If(5 > 3, "greater", "less or equal")
	fmt.Println("5 is", result)
}
```

---

## JSON Utilities

These functions provide a convenient and faster alternative to the standard library's JSON package by using [goccy/go-json](https://github.com/goccy/go-json).

### Functions

- **Unmarshal:**  
  `Unmarshal[T any](rc io.Reader) (*T, error)`  
  Reads all data from an `io.Reader`, unmarshals it into a variable of type `T`, and returns a pointer to the result.

- **UnmarshalPointer:**  
  `UnmarshalPointer[T any](in *T, rc io.Reader) (*T, error)`  
  Works similarly to `Unmarshal`, but reuses the passed pointer for potentially improved performance.

### Example

```go
package main

import (
	"fmt"
	"strings"

	"github.com/41Baloo/btils"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	jsonStr := `{"name": "Alice", "age": 30}`

	// Using Unmarshal
	person, err := btils.Unmarshal[Person](strings.NewReader(jsonStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshaled Person: %+v\n", person)

	// Using UnmarshalPointer
	var person2 Person
	_, err = btils.UnmarshalPointer(&person2, strings.NewReader(jsonStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshaled Person (pointer): %+v\n", person2)
}
```
