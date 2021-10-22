package main

import (
	"fmt"
	"github.com/Hamster601/go_order"
)

func main() {
	res:= go_order.New().
		Add("func1",nil,f1).
		Add("func2",[]string{"func1"},f2).
		Add("func3",[]string{"func2"},f3).
		Add("func4",[]string{"func3"},f4)
	res1,err := res.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res1)
}

func f1(r map[string]interface{}) (interface{},error){
	fmt.Println("test1")
	return 1,nil
}

func f2(r map[string]interface{}) (interface{},error){
	fmt.Println("test2")
	return 2,nil
}

func f3(r map[string]interface{}) (interface{},error){
	fmt.Println("test3")
	return 3,nil
}

func f4(r map[string]interface{}) (interface{},error){
	fmt.Println("test4")
	return 4,nil
}