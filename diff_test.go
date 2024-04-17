package diff

import (
	"log"
	"testing"
)

func TestUnmarshalNoDiff(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	bytes := []byte(`{"name":"merlin","age":30}`)
	var person Person
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.Name != "merlin" {
		t.Errorf("expected name to be merlin, got %s", person.Name)
	}

	if person.Age != 30 {
		t.Errorf("expected age to be 30, got %d", person.Age)
	}
}

func TestUnmarshalWithDiff(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	bytes := []byte(`{"name":"merlin", "age": 30}`)
	var person Person
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.Name != "merlin" {
		t.Errorf("expected name to be merlin, got %s", person.Name)
	}
	if person.Age != 30 {
		t.Errorf("expected age to be 30, got %d", person.Age)
	}
}

func TestUnmarshalWithLessFields(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	bytes := []byte(`{"name":"merlin"}`)
	var person Person
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.Name != "merlin" {
		t.Errorf("expected name to be merlin, got %s", person.Name)
	}
	if person.Age != 0 {
		t.Errorf("expected age to be 0, got %d", person.Age)
	}
}

func TestUnmarshalWithNoJsonTag(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	var person Person
	bytes := []byte(`{"Name":"merlin","Age":30}`)
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.Name != "merlin" {
		t.Errorf("expected name to be merlin, got %s", person.Name)
	}

	if person.Age != 30 {
		t.Errorf("expected age to be 30, got %d", person.Age)
	}
}

func TestUnmarshalWithNoExports(t *testing.T) {
	type Person struct {
		name string
		age  int
	}

	var person Person
	bytes := []byte(`{"name":"merlin","age":30}`)
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.name != "" {
		t.Errorf("expected name to be empty, got %s", person.name)
	}

	if person.age != 0 {
		t.Errorf("expected age to be 0, got %d", person.age)
	}

}

func TestUnmarshalWithEmbeddedFields(t *testing.T) {
	type Pet struct {
		Name string `json:"name"`
	}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Pets []Pet  `json:"pets"`
	}

	var person Person
	bytes := []byte(`{
    "name": "merlin",
    "age": 30,
    "pets": [
        {
            "name": "dog"
        },
        {
            "name": "cat"
        }
    ]
}`)
	diff, err := Unmarshal(bytes, &person)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(person)

	if len(diff) != 0 {
		t.Errorf("expected diff to be empty, got %v", diff)
	}

	if person.Name != "merlin" {
		t.Errorf("expected name to be merlin, got %s", person.Name)
	}

	if person.Age != 30 {
		t.Errorf("expected age to be 30, got %d", person.Age)
	}

	if len(person.Pets) != 2 {
		t.Errorf("%v", person)
		t.Errorf("expected to have 2 pets, got %d", len(person.Pets))
	}

	if person.Pets[0].Name != "dog" {
		t.Errorf("expected pet type to be dog, got %s", person.Pets[0].Name)
	}

	// check for first pet

}
