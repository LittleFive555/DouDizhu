package define

type DStrings struct {
    ID string
    Value string
}

func (data DStrings) GetID() string {
    return data.ID
}
type DStringsList struct {
    Content []DStrings
}
