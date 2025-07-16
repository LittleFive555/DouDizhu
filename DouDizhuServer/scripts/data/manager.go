package data

import (
	"DouDizhuServer/scripts/data/define"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type DataManager struct {
	Data map[string]interface{}
}

// 需要从datas文件夹读取json文件
func GetConfigStrKey[T define.DBaseData[string]](index string) T {
	var result T
	// 获取类型名称用于构建文件路径
	typeName := reflect.TypeOf((*T)(nil)).Elem().Name()
	typeName = typeName[1:]
	filePath := fmt.Sprintf("F:\\DouDizhu\\DouDizhuServer\\datas\\%s", typeName)
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err) // TODO 不许panic
	}
	typeaaa := reflect.TypeOf((*T)(nil)).Elem()
	listType := define.GetListType(typeaaa)
	listInstance := reflect.New(listType).Interface()
	err = json.Unmarshal(jsonData, listInstance)
	if err != nil {
		panic(err) // TODO 不许panic
	}
	content := reflect.ValueOf(listInstance).Elem().FieldByName("Content")
	for _, item := range content.Interface().([]T) {
		if item.GetID() == index {
			result = item
			break
		}
	}

	return result
}
