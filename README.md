# go-cross

A go cross platform library with abstractions for standard lib global functions.

## using

```bash
go get github.com/patrickhuber/go-crosos
```

### os

```go
import(
  "github.com/patrickhuber/go-cross/os"
)
func main(){
  o := os.New()
  fmt.Println(o.Executable())
}
```

### env

```go
import(
  "github.com/patrickhuber/go-cross/env"
)
func main(){
  e := env.New()
  e.Set("MY_ENV_VAR", "test")
  fmt.Println(e.Get("MY_ENV_VAR"))
}
```

```
test
```

### console

```go
import(
  "github.com/patrickhuber/go-cross/console"
)  
func main(){
  c := console.New()
  fmt.Fprintln(c.Out(), "hello world")
}
```

```
hello world
```

### filepath

```go
import(
  "github.com/patrickhuber/go-cross/filepath"
) 

func main(){
  // creates the default platform parser
  parser := filepath.NewParser()
  fp, err := parser.Parse("/some/path/to/parse")
  if err != nil{
    fmt.Fprintf("%w", err)
    os.Exit(1)
  }
  for _, seg := range fp.Segments{
    fmt.Println(seg)
  }
}
```

```
some
path
to
parse
```