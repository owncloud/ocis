package channels

import (
	"net/mail"
	"testing"
)

func Test_appendSender(t *testing.T) {
	type args struct {
		sender string
		a      mail.Address
	}

	a1, err := mail.ParseAddress("ownCloud <noreply@example.com>")
	if err != nil {
		t.Error(err)
	}
	a2, err := mail.ParseAddress("noreply@example.com")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name   string
		sender string
		want1  string
		want2  string
	}{
		{
			name:   "empty sender",
			sender: "",
			want1:  `"ownCloud" <noreply@example.com>`,
			want2:  `<noreply@example.com>`,
		},
		{
			name:   "not empty sender",
			sender: `Joe Q. Public`,
			want1:  `"Joe Q. Public via ownCloud" <noreply@example.com>`,
			want2:  `"Joe Q. Public via" <noreply@example.com>`,
		},
		{
			name:   "sender whit comma and semicolon",
			sender: `Joe, Q; Public:`,
			want1:  `"Joe, Q; Public: via ownCloud" <noreply@example.com>`,
			want2:  `"Joe, Q; Public: via" <noreply@example.com>`,
		},
		{
			name:   "sender with quotes",
			sender: `Joe Q. "Public"`,
			want1:  `"Joe Q. \"Public\" via ownCloud" <noreply@example.com>`,
			want2:  `"Joe Q. \"Public\" via" <noreply@example.com>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendSender(tt.sender, *a1); got != tt.want1 {
				t.Errorf("appendSender() = %v, want %v", got, tt.want1)
			}
			if got := appendSender(tt.sender, *a2); got != tt.want2 {
				t.Errorf("appendSender() = %v, want %v", got, tt.want2)
			}
		})
	}
}
