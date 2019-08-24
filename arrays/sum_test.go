package arrays

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"collection of 5 numbers", args{[]int{1, 2, 3, 4, 5}}, 15},
		{"collection of any size", args{[]int{1, 2, 3}}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.numbers); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumAll(t *testing.T) {
	type args struct {
		numbersToSum [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"sum all", args{[][]int{{1, 2}, {0, 9}}}, []int{3, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumAll(tt.args.numbersToSum...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumAllTails(t *testing.T) {
	type args struct {
		numbersToSum [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"sum all tails", args{[][]int{{1, 2}, {0, 9}}}, []int{2, 9}},
		{"safely sum empty slices", args{[][]int{{}, {3, 4, 5, 6}}}, []int{0, 15}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumAllTails(tt.args.numbersToSum...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumAllTails() = %v, want %v", got, tt.want)
			}
		})
	}
}
