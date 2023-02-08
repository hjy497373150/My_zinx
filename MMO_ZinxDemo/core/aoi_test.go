package core

import (
	"fmt"
	"testing"
)

func TestNewAoiManager(t *testing.T) {
	// 初始化AOIManager
	AOIManager := NewAOIManager(0,250,5,0,250,5)

	fmt.Println(AOIManager)
}