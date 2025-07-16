package define

type DataIndex interface {
    string | int
}

type DBaseData[T DataIndex] interface {
    GetID() T
}
