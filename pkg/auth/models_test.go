package auth

import "testing"

func TestPayload_ImportPermissionList(t *testing.T) {
	type args struct {
		permissions []string
		scope       string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "Permission importing",
			args: args{
				permissions: []string{"a", "b", "c", CreateToken},
				scope:       "ascope",
			},
			want: map[string][]string{
				"a":         {"ascope"},
				"b":         {"ascope"},
				"c":         {"ascope"},
				CreateToken: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pld := &Payload{}
			pld.ImportPermissionList(tt.args.permissions, tt.args.scope)
			for key, want := range tt.want {
				got, ok := pld.Permissions[key]
				if !ok || (got != nil) != (want != nil) || got != nil && got[0] != want[0] {
					t.Errorf("Payload.ImportPermisisonList error, on %v got %v, want %v", key, want, got)
					return
				}
			}
		})
	}
}
