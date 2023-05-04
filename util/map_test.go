package util

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	t.Log("testMap in")

	testIDs := NewMap()
	testIDs.Set("1", true)
	testIDs.Set("2", true)
	testIDs.Set("3", true)
	testIDs.Set("5", true)

	testIDs.UnsafeRange(func(i1, i2 interface{}) {
		fmt.Println("testID=", i1, ",val=", i2)
		t.Log("testID=", i1, ",val=", i2)
	})

	t.Log("testMap end")

	// sl := make([]jsonData, 0)
	// sl = append(sl, jsonData{
	// 	Key: "key1",
	// 	Val: "val1",
	// })
	// sl = append(sl, jsonData{
	// 	Key: "key2",
	// 	Val: "val2",
	// })
	// test := &Trust{
	// 	Name: "test1",
	// 	Data: sl,
	// }

	// res, err := json.Marshal(test)
	// if err != nil {
	// 	fmt.Println("err=", err.Error())
	// }
	// fmt.Println("res=", string(res))
}

type Trust struct {
	Name string     `json:"name"`
	Data []jsonData `json:"data"`
}

type jsonData struct {
	Key string `json:"key"`
	Val string `json:"val"`
}
