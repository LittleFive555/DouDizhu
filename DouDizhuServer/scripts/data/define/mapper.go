package define

import "reflect"

var itemToList = make(map[reflect.Type]reflect.Type)

func InitMapper() {
    registerItemToList(DConst{}, DConstList{})
    registerItemToList(DStrings{}, DStringsList{})
}

func registerItemToList(item, list interface{}) {
    itemToList[reflect.TypeOf(item)] = reflect.TypeOf(list)
}

func GetListType(item reflect.Type) reflect.Type {
    return itemToList[item]
}
