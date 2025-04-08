# monkey
monkey is a library to make patch for unit tests, this repository modified some from [gomonkey](https://github.com/agiledragon/gomonkey).\
it can patch exported/unexported methods from exported/unexported structure.

## Example
### Patch function
```go
patch := func(a ...interface{}) (int, error) {
    return fmt.Print("what?!")
}
pg := monkey.Patch(fmt.Println, patch)
defer pg.Unpatch()

// output: what?!
fmt.Println("hello!")
```

### Patch method
#### with receiver
```go
var r *bytes.Reader
patch := func(*bytes.Reader, b []byte) (int, error) {
    return 0, nil
}
pg := monkey.PatchMethod(r, "Read", patch)
defer pg.Unpatch()

reader := bytes.NewReader([]byte("hello"))
buf := make([]byte, 1024)

// output: 0 <nil>
fmt.Println(reader.Read(buf))
```
#### ignore receiver
```go
var r *bytes.Reader
patch := func(b []byte) (int, error) {
    return 0, nil
}
pg := monkey.PatchMethod(r, "Read", patch)
defer pg.Unpatch()

reader := bytes.NewReader([]byte("hello"))
buf := make([]byte, 1024)

// output: 0 <nil>
fmt.Println(reader.Read(buf))
```

## Original
https://github.com/bouk/monkey  
https://github.com/agiledragon/gomonkey
