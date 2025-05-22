package delobdriver

import (
	"testing"
)

func Test_IfCanBuildQueryWithAllComponentsWhenNotUsingAnyParameters(t *testing.T) {
	delobCtx, _ := NewContext("")

	queryBuilder := delobCtx.Select().From(Players)

	if queryBuilder.query != "SELECT * FROM Players" {
		t.Errorf("Wrong query.")
	}
}

func Test_IfCanBuildQueryWithAllComponentsSpecifiedManually(t *testing.T) {
	delobCtx, _ := NewContext("")

	queryBuilder := delobCtx.Select(Key, Elo, Events, Matches).From(Players)

	if queryBuilder.query != "SELECT Key, Elo, Events, Matches FROM Players" {
		t.Errorf("Wrong query.")
	}
}

func Test_IfCanBuildQueryWithKeysOnlyOrderByKeyDesc(t *testing.T) {
	delobCtx, _ := NewContext("")

	queryBuilder := delobCtx.Select(Key).From(Players).OrderBy(Key, Descending)

	if queryBuilder.query != "SELECT Key FROM Players ORDER BY Key DESC" {
		t.Errorf("Wrong query.")
	}
}

func Test_IfCannotOrderByEventsOrMatches(t *testing.T) {
	delobCtx, _ := NewContext("")

	queryBuilder_1 := delobCtx.Select(Key).From(Players).OrderBy(Events, Descending)

	if queryBuilder_1.query != "SELECT Key FROM Players" {
		t.Errorf("Wrong query.")
	}

	queryBuilder_2 := delobCtx.Select(Key).From(Players).OrderBy(Matches, Descending)

	if queryBuilder_2.query != "SELECT Key FROM Players" {
		t.Errorf("Wrong query.")
	}
}
