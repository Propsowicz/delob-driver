package delobdriver

import (
	"fmt"
	"testing"
)

func Test_TestCase_1(t *testing.T) {
	context, err := NewContext("Server=localhost;Port=5678;Uid=myUsername;Pwd=myPassword;")
	if err != nil {
		t.Errorf("Should be able to create delob context")
	}

	result, errFetch := context.GetPlayersOrderBy(Elo, Descending)

	if errFetch == nil {
		fmt.Println(result)
		t.Errorf("There should be no error - %s.", errFetch.Error())
	}
	if errFetch != nil {
		fmt.Println(result)
		t.Errorf("There should be no error - %s.", errFetch.Error())
	}

	if len(result) != 0 {
		fmt.Println(result)
	}
}
