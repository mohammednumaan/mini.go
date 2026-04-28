package main 

import (
	"fmt"
	"github.com/mohammednumaan/mini.go/go-cache/internal/cache"
)

func main() {
	c := cache.NewCache() 
	fmt.Println(c)

	// just a bunch of examples to verify 
	// if the methods work correctly (ideally, there should be unit tests)
	setItem := c.Set("key1", "value1")
	fmt.Println(setItem)

	getItem := c.Get("key1")
	fmt.Println(getItem)

	deleteItem := c.Delete("key1")
	fmt.Println(deleteItem)

	getDeletedItem := c.Get("key1")
	fmt.Println(getDeletedItem)
}
