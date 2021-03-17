package kits

import (
	"reflect"
	"testing"
)

func TestIsInSlice(t *testing.T) {
	type args struct {
		slice []string
		key   string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "in slice", args: args{slice: []string{"key", "value"}, key: "key"}, want: true},
		{name: "not in slice", args: args{slice: []string{"key", "value"}, key: "key2"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInSlice(tt.args.slice, tt.args.key); got != tt.want {
				t.Errorf("IsInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInSliceInt(t *testing.T) {

	type args struct {
		slice []int
		key   int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "in slice", args: args{slice: []int{1, 22, 3, 44, 5}, key: 3}, want: true},
		{name: "not in slice", args: args{slice: []int{1, 22, 3, 44, 55}, key: 66}, want: false},
		{name: "empty slice", args: args{slice: []int{}, key: 11}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInSliceInt(tt.args.slice, tt.args.key); got != tt.want {
				t.Errorf("IsInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIntersect(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{name: "slice", args: args{slice1: []string{"key"}, slice2: []string{"key", "value"}}, want: []string{"key"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceIntersect(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIntersectInSlices(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "false", args: args{slice1: []string{"key"}, slice2: []string{"value"}}, want: false},
		{name: "true", args: args{slice1: []string{"key"}, slice2: []string{"key"}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIntersectInSlices(tt.args.slice1, tt.args.slice2); got != tt.want {
				t.Errorf("IsIntersectInSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateMD5ForString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test", args: args{str: "test"}, want: "098f6bcd4621d373cade4e832627b4f6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMD5ForString(tt.args.str); got != tt.want {
				t.Errorf("GenerateMD5ForString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateMD5(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test", args: args{data: []byte("test")}, want: "098f6bcd4621d373cade4e832627b4f6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMD5(tt.args.data); got != tt.want {
				t.Errorf("GenarateMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
