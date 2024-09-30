package model

import (
	"reflect"
	"testing"
)

func Test_at(t *testing.T) {
	type args struct {
		c [][]uint
		x int
		y int
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "1x1 index",
			args: args{
				c: [][]uint{{3}},
				x: 0,
				y: 0,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "2x2 index",
			args: args{
				c: [][]uint{{1, 2}, {3, 4}},
				x: 0,
				y: 1,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "off-grid index",
			args: args{
				c: [][]uint{{1, 2}, {3, 4}},
				x: 3,
				y: 3,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := at(tt.args.c, tt.args.x, tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("at() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("at() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAutomaton(t *testing.T) {
	type args struct {
		createTransitionSet func() TransitionSet
		colouring           []Rgb
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				createTransitionSet: func() TransitionSet {
					t := NewTransitionSet()
					t.AddTransition(0, 1, func(cell Cell) bool { return true })
					t.AddTransition(1, 0, func(cell Cell) bool { return true })
					return t
				},
				colouring: []Rgb{{R: 0, G: 0, B: 0}, {R: 1, G: 1, B: 1}},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				createTransitionSet: func() TransitionSet {
					t := NewTransitionSet()
					return t
				},
				colouring: []Rgb{},
			},
			wantErr: true,
		},
		{
			name: "broken",
			args: args{
				createTransitionSet: func() TransitionSet {
					t := NewTransitionSet()
					t.AddTransition(0, 1, func(cell Cell) bool { return true })
					t.AddTransition(1, 0, func(cell Cell) bool { return true })
					return t
				},
				colouring: []Rgb{{R: 0, G: 0, B: 0}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAutomaton(tt.args.createTransitionSet(), tt.args.colouring)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAutomaton() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAutomaton_CountStates(t *testing.T) {
	type args struct {
		createTransitionSet func() TransitionSet
		colouring           []Rgb
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "2",
			args: args{
				createTransitionSet: func() TransitionSet {
					t := NewTransitionSet()
					t.AddTransition(0, 1, func(cell Cell) bool {
						return true
					})
					return t
				},
				colouring: []Rgb{{R: 0, G: 0, B: 0}, {R: 0, G: 0, B: 0}},
			},
			want: 2,
		},
		{
			name: "3",
			args: args{
				createTransitionSet: func() TransitionSet {
					t := NewTransitionSet()
					t.AddTransition(0, 1, func(cell Cell) bool {
						return true
					})
					t.AddTransition(0, 2, func(cell Cell) bool {
						return true
					})
					return t
				},
				colouring: []Rgb{{R: 0, G: 0, B: 0}, {R: 0, G: 0, B: 0}, {R: 0, G: 0, B: 0}},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAutomaton(tt.args.createTransitionSet(), tt.args.colouring)
			if err != nil {
				t.Errorf("Error during creation of test automaton = %v", err)
				return
			}
			if got := a.CountStates(); got != tt.want {
				t.Errorf("Automaton.CountStates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutomaton_Step(t *testing.T) {
	type args struct {
		c               [][]uint
		createAutomaton func() *Automaton
	}
	tests := []struct {
		name string
		args args
		want [][]uint
	}{
		{
			name: "checkerboard",
			args: args{
				c: [][]uint{{1, 0}, {0, 1}},
				createAutomaton: func() *Automaton {
					t := NewTransitionSet()
					t.AddTransition(0, 1, func(cell Cell) bool { return true })
					t.AddTransition(1, 0, func(cell Cell) bool { return true })

					c := []Rgb{{R: 0, G: 0, B: 0}, {R: 1, G: 1, B: 1}}

					a, _ := NewAutomaton(t, c)
					return a
				},
			},
			want: [][]uint{{0, 1}, {1, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args.createAutomaton()
			if got := a.Step(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Automaton.Step() = %v, want %v", got, tt.want)
			}
		})
	}
}
