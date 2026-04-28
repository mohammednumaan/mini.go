package main 

import (
	"fmt"
	"time"
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

	// now i check if the item is "cleaned up" after the TTL expires
	// for now the TTL is set to 1 minute
	c.Set("key2", "value2")
	fmt.Println(c.Get("key2"))

	c.CleanUp()
	fmt.Println("Before 5 seconds", c.Get("key2"))

	time.Sleep(5 * time.Second)
	c.CleanUp()
	fmt.Println("After 5 seconds", c.Get("key2"))
}
