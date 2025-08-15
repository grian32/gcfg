# GCFG

GCFG is a config format that aims to be simple & readable. It draws some inspiration from TOML, JSON & a couple of custom formats.

GCFG is designed for configuration, not for general purpose data storage, and as such it doesn't support arbitrary nesting.

## Features

GCFG supports the following value types: integers, floats, strings, booleans, arrays, pairs, and nil.

### Pairs

Pairs are a tuple of two values of the same type. 
```gcfg
pair = (3, 3)
```

The GCFG Library also provides a Pair type.

### Sections

Sections are named blocks that group values together. They cannot contain other sections.

```gcfg
Block { 
    a = 1
}
```

### Arrays

Arrays require all elements to be of the same type, and there is specific syntax for arrays of sections:
```
[SecArr] { 
    a = 1
}
    
[SecArr] { 
    a = 2
}
```

The above will create an array named `SecArr` containing anonymous sections with the same structure.

## Examples

### Free Standing Configs

```gcfg
foo = "true"
bar = 3
baz = false
faz = 4.4
aaa = ["a", "b", "c"]
bbb = (3.3, 3.3)
```

This maps to:
```go
type Config struct { 
    Foo string `gcfg:"foo"`
    Bar int32 `gcfg:"bar"`
    Baz bool `gcfg:"baz"`
    Faz float32 `gcfg:"faz"`
    Aaa []string `gcfg:"aaa"`
    Bbb Pair[float32, float32] `gcfg:"bbb"`
}
```

### Sections

```gcfg
Baz { 
    foo = -3
    faz = 4
}

Bar { 
    aaa = -3
    bbb = 4
}
```

This maps to:

```go
type Config struct {
    Baz Baz `gcfg:"Baz"`
    Bar Bar `gcfg:"Bar"`
}

type Baz struct {
    Foo int32 `gcfg:"foo"`
    Faz int16 `gcfg:"faz"`
}

type Bar struct {
    Aaa int32 `gcfg:"aaa"`
    Bbb int16 `gcfg:"bbb"`
}
```

### Arrays of Sections

```gcfg
[Foo] {
    bar = 5
}

[Foo] {
    bar = 6
}

[Foo] {
    bar = 7
}

[Foo] {
    bar = 8
}
```

This maps to:
```go
type Config struct {
    Foo []Foo `gcfg:"Foo"`
}

type Foo struct {
    Bar uint16 `gcfg:"bar"`
}
```