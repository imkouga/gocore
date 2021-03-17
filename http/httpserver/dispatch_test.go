package httpserver

import (
	"testing"
)

func Test_compilePrefixPathRegexp(t *testing.T) {

	prefix := "/hello"
	reg, err := compilePrefixPathRegexp(prefix)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		template string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{"hello"}, wantErr: false},
		{name: "test2", args: args{"/Hello"}, wantErr: false},
		{name: "test3", args: args{"/helo"}, wantErr: false},
		{name: "test4", args: args{"/hello"}, wantErr: false},
		{name: "test5", args: args{"/hello/"}, wantErr: true},
		{name: "test6", args: args{"/hello/xaw3"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reg.MatchString(tt.args.template)
			if tt.wantErr != got {
				t.Errorf("%s compilePrefixPathRegexp()  wantErr %v, got %v", tt.name, tt.wantErr, got)
				return
			}

		})
	}
}
