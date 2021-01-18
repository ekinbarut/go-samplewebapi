package main

import "testing"

func TestReadJsonFile(t *testing.T) {
	objects := readJsonFile()
	obj := objects[1]

	expectedId := "2"
	expectedTitle := "Hello 2"

	if expectedId != obj.Id {
		t.Errorf("Object.Id == %q, want %q",
			obj.Id, expectedId)
	}

	if expectedTitle != obj.Title {
		t.Errorf("Object.Title == %q, want %q",
			obj.Title, expectedTitle)
	}

}
