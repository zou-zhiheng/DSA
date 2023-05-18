package Goroutine

import (
	"fmt"
	"reflect"
	"sync"
)

type Obj struct {
	Name string //姓名
	Sex  string //性别
	Age  int    //年龄
}

func (r *Obj) Get(key string) interface{} {

	retH:=reflect.TypeOf(*r)
	//获取接口体里面的名称
	for i:=0;i<retH.NumField();i++{
		field:=retH.Field(i)
		if field.Name==key{
			valStu:=reflect.ValueOf(*r)
			return valStu.FieldByName(field.Name)
		}
	}

	return nil
}


var ObjPool=sync.Pool{
	New: func() interface{} {
		//创建新的实例
		return &Obj{}
	},
}

func Demo(){
	test:=Obj{
		Name: "test",
		Sex: "男",
		Age: 20,
	}

	t:=test.Get("Age")
	fmt.Println(t)
}

