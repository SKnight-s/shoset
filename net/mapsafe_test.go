package net_test

import (
	"fmt"
	"shoset/net"
	"strconv"
	"testing"
	"time"
)

// TestMapSafeCRUD : test MapSafe crud functions
func TestMapSafeCRUD(t *testing.T) {
	m := net.NewMapSafe()
	m.Set("a", 23).Set("b", 43).Set("c", 11)
	fmt.Printf("Initial state : ")
	m.Iterate(
		func(key string, val interface{}) {
			fmt.Printf(" - %s: %d\n", key, val.(int))
		},
	)
	m.Set("a", 454433)
	fmt.Printf("After update : ")
	m.Iterate(
		func(key string, val interface{}) {
			fmt.Printf(" - %s: %d\n", key, val.(int))
		},
	)
	m.Delete("b")
	fmt.Printf("After delete : ")
	m.Iterate(
		func(key string, val interface{}) {
			fmt.Printf(" - %s: %d\n", key, val.(int))
		},
	)
}

// TestFold : test MapSafe fold function
func TestFold(t *testing.T) {
	m := net.NewMapSafe()
	m.Set("a", 23).Set("b", 43).Set("c", 11)
	strValue := m.Fold(
		func(key string, val interface{}, str interface{}) interface{} {
			str = fmt.Sprintf("%s, %s", str, val)
			fmt.Printf(" - %s: %d\n", key, val.(int))
			return str
		},
		"")
	fmt.Printf("strValue : %s\n", strValue)
}

// TestConcurrency
func TestConcurrency(t *testing.T) {
	m := net.NewMapSafe()
	for i := 0; i < 10; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	fmt.Printf("test Concurrency\n")
	go m.Iterate(
		func(key string, val interface{}) {
			time.Sleep(time.Millisecond * time.Duration(10))
			fmt.Printf("%s, %d\n", key, val.(int))
		})
	time.Sleep(time.Millisecond * time.Duration(20))
	fmt.Printf("after Iterate\n")
	m.Set("a", 101)
	m.Set("b", 102)
	fmt.Printf("after Set\n")
	m.Iterate(
		func(key string, val interface{}) {
			fmt.Printf("%s, %d\n", key, val.(int))
		})

}
