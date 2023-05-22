package auth

import "testing"

// TestAccount_GetID - test for Account.GetID()
func TestAccount_GetID(t *testing.T) {
	type fields struct {
		ID       string
		Login    string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "1",
				Login:    "login",
				Password: "password",
			},
			want: "1",
		},
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "2",
				Login:    "login",
				Password: "password",
			},
			want: "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				ID:           tt.fields.ID,
				Login:        tt.fields.Login,
				PasswordHash: tt.fields.Password,
			}
			if got := a.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})

	}
}

// TestAccount_GetLogin - test for Account.GetLogin()
func TestAccount_GetLogin(t *testing.T) {
	type fields struct {
		ID       string
		Login    string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "1",
				Login:    "login1",
				Password: "password",
			},
			want: "login1",
		},
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "2",
				Login:    "login2",
				Password: "password",
			},
			want: "login2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				ID:           tt.fields.ID,
				Login:        tt.fields.Login,
				PasswordHash: tt.fields.Password,
			}
			if got := a.GetLogin(); got != tt.want {
				t.Errorf("GetLogin() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestAccount_SetPasswordHash(t *testing.T) {
	type fields struct {
		ID       string
		Login    string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "1",
				Login:    "login1",
				Password: "password1",
			},
			want: "password1",
		},
		{
			name: "TestAccount_GetID",
			fields: fields{
				ID:       "2",
				Login:    "login2",
				Password: "password2",
			},
			want: "password2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				ID:    tt.fields.ID,
				Login: tt.fields.Login,
			}
			a.SetPasswordHash(tt.fields.Password)
			if got := a.PasswordHash; got != tt.want {
				t.Errorf("GetLogin() = %v, want %v", got, tt.want)
			}
		})

	}
}
