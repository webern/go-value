// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"testing"

	"github.com/webern/tcore"
)

func TestNewObjectValue(t *testing.T) {
	v := NewObjectValue(nil)

	gotS := v.Type().String()
	wantS := ObjectType.String()
	if msg, ok := tcore.TAssertString("", gotS, wantS); !ok {
		t.Error(msg)
	}

	o := v.Object()
	if o == nil {
		t.Error("value.Object() should never return nil, but did")
	}

	v = Value{}
	o = v.Object()
	if o == nil {
		t.Error("value.Object() should never return nil, but did")
	}
}
