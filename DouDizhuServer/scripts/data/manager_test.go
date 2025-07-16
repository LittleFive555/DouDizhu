package data

import (
	"DouDizhuServer/scripts/data/define"
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	define.InitMapper()
	consts := GetConfigStrKey[define.DConst]("AccountMinLength")
	fmt.Println(consts.Value)
}
