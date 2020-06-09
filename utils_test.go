package ddns

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "main", want: "175.15.137.206", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
