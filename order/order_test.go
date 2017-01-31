package main

import (
	"fmt"
	"testing"
)

func TestOrder_String(t *testing.T) {
	cases := []struct {
		in   Order
		want string
	}{
		{
			in: Order{
				Type:      Market,
				Side:      BuySide,
				Volume:    1,
				Remaining: 1,
			},
			want: "market order to buy 1 units at market price, with 1 remaining",
		},
		{
			in: Order{
				id:        2,
				Type:      Market,
				Side:      SellSide,
				Volume:    1,
				Remaining: 0,
			},
			want: "[id 2] market order to sell 1 units at market price, with 0 remaining",
		},
		{
			in: Order{
				id:        3,
				cancelled: true,
				Type:      Limit,
				Side:      SellSide,
				Volume:    3,
				Remaining: 2,
				Limit:     3.4,
			},
			want: "[id 3] [cancelled] limit order to sell 3 units >= $3.40, with 2 remaining",
		},
		{
			in: Order{
				id:            4,
				stopTriggered: true,
				executed:      true,
				Type:          Stop,
				Side:          BuySide,
				Volume:        3,
				Remaining:     0,
				Limit:         3.4,
			},
			want: "[id 4] [executed] [triggered] stop order to buy 3 units if price goes > 3.40, with 0 remaining",
		},
		{
			in: Order{
				id:       43,
				Type:     Cancel,
				ToCancel: 42,
			},
			want: "[id 43] cancel order that disables #42",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			got := tc.in.String()
			if got != tc.want {
				t.Errorf("got %q; want %q", got, tc.want)
			}
		})
	}
}

func BenchmarkExecute(b *testing.B) {
	benchmarks := []struct {
		name    string
		float   float64
		fmt     byte
		prec    int
		bitSize int
	}{
		{"Decimal", 33909, 'g', -1, 64},
		{"Float", 339.7784, 'g', -1, 64},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// TODO: Benchmark order execution:
				// AppendFloat(dst[:0], bm.float, bm.fmt, bm.prec, bm.bitSize)
			}
		})
	}
}
