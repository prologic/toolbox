package mo_url

import "testing"

func TestNewUrl(t *testing.T) {
	u, err := NewUrl("https://www.dropbox.com")
	if err != nil {
		t.Error(err)
	} else if u.String() != "https://www.dropbox.com" {
		t.Error(u.String())
	}
}
