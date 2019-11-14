package domain

import "testing"

func TestGeneratingAuthorizationCode(t *testing.T) {
	builder := &authorizationBuilder{
		responseType: "code",
		clientId:     "xaaaaa",
	}
	az, err := builder.Build()
	if err != nil {
		t.Errorf("got an error %v but didn't want one", err.Error())
	}

	if az.AuthzCode == "" {
		t.Errorf("wanted a code but didn't get one")
	}

	//if az.state != "" && hoge.state == "" {
	//	t.Errorf("wanted a state but didn't get one")
	//}
}
