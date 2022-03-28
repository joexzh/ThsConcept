package db

import "testing"

func TestParamList(t *testing.T) {
	t.Log("TestParamList")
	list := []int{1, 2, 3}

	expectedSql := "(?,?,?)"
	expectedParams := []interface{}{1, 2, 3}

	listSql, params := ParamList(list...)

	if listSql != expectedSql {
		t.Errorf("Expected sql %s, got %s", expectedSql, listSql)
	}

	if len(params) != len(expectedParams) {
		t.Errorf("Expected %d params, got %d", len(expectedParams), len(params))
	}
}
