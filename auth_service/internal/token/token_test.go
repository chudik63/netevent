package token

// func TestNewToken(t *testing.T) {
// 	type Data struct {
// 		name  string
// 		token string
// 	}
// 	data := []Data{
// 		Data{
// 			name: "useraфвыаввв",
// 			//token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.cThIIoDvwdueQB468K5xDc5633seEFoqwxjF_xSJyQQ`,
// 		},
// 	}
// 	for _, v := range data {
// 		toc, err := NewToken(v.name)
// 		if err != nil {
// 			t.Errorf("error %v", err)
// 		}
// 		time.Sleep(time.Second)
// 		if res, err := ValidToken(toc); err != nil || !res {
// 			t.Errorf("error %v %v", err, res)
// 		}
// 	}
// }

// func TestGetName(t *testing.T) {
// 	data := []int64{
// 		1, 2, 3, 4, 5,
// 	}
// 	for _, v := range data {
// 		tock, err := NewTokens(v, "user")
// 		if err != nil {
// 			t.Error(err.Error())
// 		}
// 		id, err := GetIdToken(tock.AccessTkn)
// 		if err != nil {
// 			t.Error(err.Error())
// 		}
// 		if id != v {
// 			t.Errorf("name not eauel %s = %s", id, v)
// 		}
// 	}
// }
