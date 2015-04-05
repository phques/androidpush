package gopush

import (
	"os"
	"testing"
)

func TestCheckBooksDir(t *testing.T) {
	// sanity check !!
	if _, err := os.Stat("./Books"); err == nil {
		t.Fatal("ooops, ./Books exists !!?? aborting test")
	}

	// look for ./Books ...
	param := &InitParam{}
	param.Documents = "./Documents"
	checkBooksDir(param)

	// .. should not be found / there
	if len(param.Books) != 0 {
		t.Errorf("found ./Books ! (should not be there, check your directory!): [%v]\n", param.Books)
	}

	// create ./Books & test again
	os.Mkdir("./Books", 0777)
	checkBooksDir(param)
	if len(param.Books) == 0 {
		t.Errorf("failed to find ./Books ! : [%v]\n", param.Books)
	}

	// delete ./Books
	os.Remove("./Books")
}
