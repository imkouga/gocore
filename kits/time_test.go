package kits

import (
	"reflect"
	"testing"
	"time"
)

func TestStatCostTimeString(t *testing.T) {

	type args struct {
		st      time.Time
		nameStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "TestStatCostTimeString",
			args: args{st: time.Now().Add(-1 * time.Second), nameStr: "TestStatCostTimeString"},
			want: "TestStatCostTimeString cost time 1.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatCostTimeString(tt.args.st, tt.args.nameStr); got != tt.want {
				t.Errorf("StatCostTimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMonthStartTimeForTime(t *testing.T) {

	type args struct {
		t time.Time
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{
			name: "2020-08-18 15:40:00",
			args: args{t: time.Date(2020, 8, 18, 15, 40, 0, 0, now.Location())},
			want: time.Date(2020, 8, 1, 0, 0, 0, 0, now.Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMonthStartTimeForTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMonthStartTimeForTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMonthEndTimeForTime(t *testing.T) {

	type args struct {
		t time.Time
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{
			name: "2020-08-18 15:40:00",
			args: args{t: time.Date(2020, 8, 18, 15, 40, 0, 0, now.Location())},
			want: time.Date(2020, 8, 31, 23, 59, 59, 0, now.Location()),
		},
		{
			name: "2020-02-18 15:40:00",
			args: args{t: time.Date(2020, 2, 18, 15, 40, 0, 0, now.Location())},
			want: time.Date(2020, 2, 29, 23, 59, 59, 0, now.Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMonthEndTimeForTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMonthEndTimeForTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
