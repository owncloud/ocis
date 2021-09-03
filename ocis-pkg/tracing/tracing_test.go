package tracing

import "testing"

func Test_parseAgentConfig(t *testing.T) {
	type args struct {
		ae string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "docker-style config",
			args: args{
				ae: "docker-jaeger:6666",
			},
			want:    "docker-jaeger",
			want1:   "6666",
			wantErr: false,
		},
		{
			name: "agent in an url config",
			args: args{
				ae: "https://example-agent.com:6666",
			},
			want:    "example-agent.com",
			want1:   "6666",
			wantErr: false,
		},
		{
			name: "agent as ipv4",
			args: args{
				ae: "127.0.0.1:6666",
			},
			want:    "127.0.0.1",
			want1:   "6666",
			wantErr: false,
		},
		{
			name: "no hostname config should error",
			args: args{
				ae: ":6666",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "no hostname nor port but separator should error",
			args: args{
				ae: ":",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseAgentConfig(tt.args.ae)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAgentConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseAgentConfig() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseAgentConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
