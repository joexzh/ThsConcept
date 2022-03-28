package db

import "testing"

func TestParamListThreeParams(t *testing.T) {
	t.Log("TestParamListThreeParams")
	list := []int{1, 2, 3}

	expectedSql := "(?,?,?)"
	expectedVals := []interface{}{1, 2, 3}

	listSql, params := ParamList(list...)

	if listSql != expectedSql {
		t.Errorf("Expected sql %s, got %s", expectedSql, listSql)
	}

	if len(params) != len(expectedVals) {
		t.Errorf("Expected %d params, got %d", len(expectedVals), len(params))
	}
}

func TestParamListZeroParams(t *testing.T) {
	t.Log("TestParamListThreeParams")
	list := []int{}

	expectedSql := "(null)"
	expectedVals := make([]interface{}, 0)

	listSql, params := ParamList(list...)

	if listSql != expectedSql {
		t.Errorf("Expected sql %s, got %s", expectedSql, listSql)
	}
	if len(params) != len(expectedVals) {
		t.Errorf("Expected %d params, got %d", len(expectedVals), len(params))
	}
}
