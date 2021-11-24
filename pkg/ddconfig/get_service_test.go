package ddconfig

import "testing"

func TestGetService(t *testing.T) {
	// Set the env var.
	t.Setenv("DD_SERVICE", "ddservice")

	tests := []struct {
		name   string
		opts   []GetServiceOption
		setEnv func()
		want   string
	}{
		{
			name: "use prefix only",
			opts: []GetServiceOption{
				WithServicePrefix("prefix-"),
			},
			want: "prefix-ddservice",
		},
		{
			name: "use suffix only",
			opts: []GetServiceOption{
				WithServiceSuffix("-suffix"),
			},
			want: "ddservice-suffix",
		},
		{
			name: "use prefix, suffix both",
			opts: []GetServiceOption{
				WithServicePrefix("prefix-"),
				WithServiceSuffix("-suffix"),
			},
			want: "prefix-ddservice-suffix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetService(tt.opts...); got != tt.want {
				t.Errorf("GetService() = %v, want %v", got, tt.want)
			}
		})
	}
}
