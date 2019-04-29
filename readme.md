go-value
========

Value wraps a fixed set of very common types in a static-typed-friendly way. It frees you up from using `interface{}` for a lot of use cases where you don't know the type you will receive (e.g. JSON).

It looks like this:

```
v := Value{}
v.IsNull()                  // true
v.IsInt()                   // false
v.SetInt(5)
v.IsInt()                   // true
v.Int()                     // 5
v.String()                  // ""
v.IsString()                // false
v.SetString("Hello Value")
v.IsString()                // true
v.String()                  // "Hello Value"
```

This is a simple package. I may break the API though so this is version zero.

#### TODO

In particular I want to change the functions like `CoerceToBool` to be pure standalone functions instead of member functions.
