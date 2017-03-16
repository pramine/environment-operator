package bitesize

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindByNameExist(t *testing.T) {
	var svc = Services{
		{Name: "ads"},
		{Name: "vpd"},
		{Name: "ooo"},
	}

	if svc.FindByName("vpd") == nil {
		t.Errorf("Expected service, got nil")
	}
}

func TestFindByNameNotFound(t *testing.T) {
	var svc = Services{
		{Name: "ads"},
		{Name: "vpd"},
		{Name: "ooo"},
	}

	if s := svc.FindByName("aaa"); s != nil {
		t.Errorf("Expected nil, got %v", s)
	}
}

func TestServiceSortInterface(t *testing.T) {
	var s = Services{
		{Name: "b"},
		{Name: "a"},
		{Name: "c"},
	}

	var expected = Services{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}

	sort.Sort(s)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Service sort invalid, got %q", s)
	}
}
