# Description
Wrapper for json.Unmarshal that returns the extra fields that weren't unmarshalled.
Useful for keeping an eye on, and catching unexpected differences in API responses instead of just ignoring them.

# Usage
```go
import "github.com/glensargent/diff"

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// this json is adding an extra field, city, that's uncaught by the struct
bytes := []byte(`{"name":"merlin", "age": 30, "city": "new york"}`)
// person will unmarshal as intended
var person Person
// diff returns a map of the fields that weren't unmarshalled in to the struct, in this case, city
difference, err := Unmarshal(bytes, &person)
if err != nil {
    t.Fatal(err)
}

log.Println("person:", person)
log.Println("extra fields:", difference)
```

# How to get 
```bash
go get github.com/glensargent/diff
```
