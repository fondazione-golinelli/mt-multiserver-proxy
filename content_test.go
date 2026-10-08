package proxy

import (
	"testing"

	"github.com/HimbeerserverDE/mt"
)

func TestInstanceNodeMappingUsesItsOwnDefinitions(t *testing.T) {
	cc := &ClientConn{globalNodeDefs: map[string]mt.Content{"shared_test:stone": 42}}
	mapped, missing, _ := cc.setServerParam0Map("instance-1", "shared", []mt.NodeDef{{Name: "test:stone", Param0: 100}})
	if mapped != 1 || missing != 0 || cc.p0Map["instance-1"][100] != 42 {
		t.Fatalf("unexpected instance mapping: mapped=%d missing=%d map=%v", mapped, missing, cc.p0Map)
	}
	if cc.p0Map["instance-1"][mt.Air] != mt.Air {
		t.Fatal("air node mapping was not preserved")
	}
}
