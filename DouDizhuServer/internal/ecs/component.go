package ecs

type Component struct {
	typeId ComponentTypeId
	Data   IComponentData
}

func NewComponent(typeId ComponentTypeId, data IComponentData) *Component {
	return &Component{
		typeId: typeId,
		Data:   data,
	}
}

type IComponentData interface {
	IsComponent()
}

type ComponentTypeId int32

type ComponentTypeIdSlice []ComponentTypeId

func (s ComponentTypeIdSlice) Len() int {
	return len(s)
}

func (s ComponentTypeIdSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s ComponentTypeIdSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
