package define

type DConst struct {
    ID string
    Value string
}

func (data DConst) GetID() string {
    return data.ID
}
type DConstList struct {
    Content []DConst
}
