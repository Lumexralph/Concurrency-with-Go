package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

type stub struct {
	toFetch          []Thing
	curr             int
	fetchErr, putErr bool
	gotPut           []Thing
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

// MaybeMove tests implementation
func (s *stub) mayBeFetch() (Thing, bool, error) {
	if s.fetchErr == true {
		return nil, false, errors.New("something went wrong fetching things")
	}

	if s.curr >= len(s.toFetch) {
		return nil, false, nil
	}

	t := s.toFetch[s.curr]
	s.curr++
	return t, true, nil
}

func (s *stub) mayBePut(t Thing) error {
	if s.putErr {
		return errors.New("could not continue putting thing")
	}

	s.gotPut = append(s.gotPut, t)
	return nil
}

func TestMaybeMove(t *testing.T) {
	t.Run("no errors during fetch and put", func(t *testing.T) {
		testCase := struct {
			stub *stub
			want []Thing
		}{
			stub: &stub{
				toFetch: []Thing{1, 2, 3},
			},
			want: []Thing{1, 2, 3},
		}

		if err := MaybeMove(testCase.stub.mayBeFetch, testCase.stub.mayBePut); err != nil {
			t.Errorf("MaybeMove() got err %v; want %v", err, nil)
		}

		if diff := cmp.Diff(testCase.want, testCase.stub.gotPut); diff != "" {
			t.Errorf("MaybeMove() mismatch (-want +got):\n%s", diff)
		}
	})

	testCases := []struct{
		name string
		stub *stub
		want []Thing
	}{
		{
			name: "errors in put",
			stub: &stub{
				// I discovered a bug when []Thing{1, 2} is supplied, the second
				// execution of fetch runs, why? It seems the second call to fetch was
				// placed on the stack and executed leading to the panic on closing
				// an already closed channel.
				toFetch: []Thing{1, 2},
				putErr:  true,
			},
		},
		{
			name: "errors in fetch",
			stub: &stub{
				toFetch: []Thing{1, 2},
				fetchErr:  true,
			},
		},
		{
			name: "errors in fetch and put",
			stub: &stub{
				toFetch: []Thing{1, 2},
				fetchErr:  true,
				putErr: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := MaybeMove(tc.stub.mayBeFetch, tc.stub.mayBePut); err == nil {
				t.Errorf("MaybeMove() should get error, got %v", err)
			}

			if diff := cmp.Diff(tc.want, tc.stub.gotPut); diff != "" {
				t.Errorf("MaybeMove() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
