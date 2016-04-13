package main

import "testing"

var cases = []struct {
	in   Data
	want Data
}{
	{
		in:   Data{0},
		want: Data{0},
	},
	{
		in:   Data{0, 1, 2},
		want: Data{0, 1, 2},
	},
	{
		in:   Data{1, 3, 2},
		want: Data{1, 2, 3},
	},
	{
		in:   Data{25, 1, 3, 2},
		want: Data{1, 2, 3, 25},
	},
}

func TestInsertion(t *testing.T) {
	for i, tt := range cases {
		d := make(Data, len(tt.in))
		copy(d, tt.in)
		d.Insertion()
		if !d.Equal(tt.want) {
			t.Errorf("[%d] %v.Insertion() => %v; want %v\n", i, tt.in, d, tt.want)
		}
	}
}

func TestMerge(t *testing.T) {
	for i, tt := range cases {
		d := make(Data, len(tt.in))
		copy(d, tt.in)
		d.Merge(0, len(d))
		if !d.Equal(tt.want) {
			t.Errorf("[%d] %v.Merge() => %v; want %v\n", i, tt.in, d, tt.want)
		}
	}
}
