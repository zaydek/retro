package expect

import (
	"reflect"
	"testing"
)

func NotDeepEqual(t *testing.T, x, y interface{}) {
	if !reflect.DeepEqual(x, y) {
		return
	}
	s1, t1 := x.(string)
	s2, t2 := y.(string)
	if t1 && t2 {
		t.Fatalf("got %q want %q", s1, s2)
	}
	t.Fatalf("got %+v want %+v", x, y)
}

func DeepEqual(t *testing.T, x, y interface{}) {
	if reflect.DeepEqual(x, y) {
		return
	}
	s1, t1 := x.(string)
	s2, t2 := y.(string)
	if t1 && t2 {
		t.Fatalf("got %q want %q", s1, s2)
	}
	t.Fatalf("got %+v want %+v", x, y)
}
