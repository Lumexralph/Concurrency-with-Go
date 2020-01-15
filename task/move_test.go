package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type stub struct {
	toFetch []Thing
	curr    int
	gotPut  []Thing
}

func (s *stub) fetch() (Thing, bool) {
	if s.curr >= len(s.toFetch) {
		return nil, false
	}
	t := s.toFetch[s.curr]
	s.curr++
	return t, true
}

func (s *stub) put(t Thing) {
	s.gotPut = append(s.gotPut, t)
}

func TestMoveCtx(t *testing.T) {
	tests := []struct {
		name string
		stub *stub
		want []Thing
	}{
		{
			name: "empty",
			stub: new(stub),
			want: []Thing{},
		},
		{
			name: "single",
			stub: &stub{
				toFetch: []Thing{1},
			},
			want: []Thing{1},
		},
		{
			name: "multiple",
			stub: &stub{
				toFetch: []Thing{1, 2, 3},
			},
			want: []Thing{1, 2, 3},
		},
	}

	for _, cancelFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("Context cancelled %t", cancelFirst), func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					var wantErr error
					var want []Thing

					ctx, cancel := context.WithCancel(context.Background())
					if cancelFirst {
						cancel()
						wantErr = context.Canceled
					} else {
						want = tt.want
					}
					defer cancel()

					if err := MoveCtx(ctx, tt.stub.fetch, tt.stub.put); err != wantErr {
						t.Errorf("MoveCtx() got error %v; want %v", err, wantErr)
					}

					got := tt.stub.gotPut

					// This is because the playground can only use reflect.DeepEqual, which
					// doesn't work properly for comparing nil and empty slices.
					if len(want) == 0 {
						if len(got) != 0 {
							t.Errorf("Move() got values put %v; want empty", got)
						}
						return
					}

					if got, want := tt.stub.gotPut, tt.want; !reflect.DeepEqual(got, want) {
						t.Errorf("MoveCtx() got values put %v; want %v", got, want)
					}
				})
			}

		})

	}
}
