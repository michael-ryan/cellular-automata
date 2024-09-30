package model

import "testing"

func TestCell_Neighbour(t *testing.T) {
	type fields struct {
		x     int
		y     int
		cells [][]uint
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "off-grid",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 2, 3},
					{4, 0, 5},
				},
			},
			args: args{
				x: 1,
				y: 0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "north",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 2, 3},
					{4, 0, 5},
					{6, 7, 8},
				},
			},
			args: args{
				x: 0,
				y: 1,
			},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cell{
				x:     tt.fields.x,
				y:     tt.fields.y,
				cells: tt.fields.cells,
			}
			got, err := c.Neighbour(tt.args.x, tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cell.Neighbour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cell.Neighbour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCell_CountNeighbours(t *testing.T) {
	type fields struct {
		x     int
		y     int
		cells [][]uint
	}
	type args struct {
		target uint
		moore  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "none moore",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
			},
			args: args{
				target: 0,
				moore:  true,
			},
			want: 0,
		},
		{
			name: "none von neumann",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
			},
			args: args{
				target: 0,
				moore:  false,
			},
			want: 0,
		},
		{
			name: "edge moore",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
				},
			},
			args: args{
				target: 1,
				moore:  true,
			},
			want: 5,
		},
		{
			name: "edge von neumann",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
				},
			},
			args: args{
				target: 1,
				moore:  false,
			},
			want: 3,
		},
		{
			name: "normal moore",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
			},
			args: args{
				target: 1,
				moore:  true,
			},
			want: 8,
		},
		{
			name: "normal von neumann",
			fields: fields{
				x: 1,
				y: 1,
				cells: [][]uint{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
			},
			args: args{
				target: 1,
				moore:  false,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cell{
				x:     tt.fields.x,
				y:     tt.fields.y,
				cells: tt.fields.cells,
			}
			if got := c.CountNeighbours(tt.args.target, tt.args.moore); got != tt.want {
				t.Errorf("Cell.CountNeighbours() = %v, want %v", got, tt.want)
			}
		})
	}
}
