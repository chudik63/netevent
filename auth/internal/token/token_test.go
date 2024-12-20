package token

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTocken(t *testing.T) {
	type Data struct {
		name   string
		tocken string
	}
	data := []Data{
		Data{
			name: "useraфвыаввв",
			//tocken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.cThIIoDvwdueQB468K5xDc5633seEFoqwxjF_xSJyQQ`,
		},
	}
	for _, v := range data {
		toc, err := NewToken(v.name)
		if err != nil {
			t.Errorf("error %v", err)
		}
		time.Sleep(time.Second)
		if res, err := ValidTocken(toc); err != nil || !res {
			t.Errorf("error %v %v", err, res)
		}
	}

	toc, _ := NewToken("user")
	fmt.Println(toc)
}

func TestGetName(t *testing.T) {
	data := []string{
		"anton", "andrei", "ismail", "arkadii", "Daniil",
	}
	for _, v := range data {
		tock, err := NewToken(v)
		if err != nil {
			t.Error(err.Error())
		}
		name, err := GetNameToken(tock)
		if err != nil {
			t.Error(err.Error())
		}
		if name != v {
			t.Errorf("name not eauel %s = %s", name, v)
		}
	}
}
