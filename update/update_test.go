package update

import "testing"

func TestUpdate_CurrentVersion(t *testing.T) {
	type fields struct {
		source  string
		target  string
		service string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				source:  "https://application-repo.oss-cn-guangzhou.aliyuncs.com/processors/",
				target:  "/usr/bin/Processor",
				service: "service",
				name:    "Processor",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &Update{
				Source:  tt.fields.source,
				Target:  tt.fields.target,
				Service: tt.fields.service,
				Name:    tt.fields.name,
			}
			if got := update.CurrentVersion(); got != tt.want {
				t.Errorf("CurrentVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate_GetRemoteVersion(t *testing.T) {
	type fields struct {
		source  string
		target  string
		service string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				source:  "https://application-repo.oss-cn-guangzhou.aliyuncs.com/processors/",
				target:  "/usr/bin/Processor",
				service: "service",
				name:    "Processor-test",
			},
			want: "D41D8CD98F00B204E9800998ECF8427E",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &Update{
				Source:  tt.fields.source,
				Target:  tt.fields.target,
				Service: tt.fields.service,
				Name:    tt.fields.name,
			}
			if got := update.GetRemoteVersion(); got != tt.want {
				t.Errorf("GetRemoteVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
