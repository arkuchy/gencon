# gencon [![Go Reference](https://pkg.go.dev/badge/github.com/ari1021/gencon.svg)](https://pkg.go.dev/github.com/ari1021/gencon)

`gencon` consists of two command, `gencon` and `fixgencon`.

- `gencon`
    - finds type constraints `any` and gives hints of type constraints.
- `fixgencon`
    - finds type constraints `any` and fix `any` to appropriate type constraints.

## requirement

`gencon` and `fixgencon` requires Go 1.18+.

## gencon command

```go

package a

func f[T any](t T) {} // want "should not use 'any'. hint: string|~int"

func invoke() { // OK
	f("gopher")
	f(1)
	type MyInt int
	f(MyInt(18))
}
```

### install

```sh
$ go install github.com/ari1021/gencon/cmd/gencon@latest
```

### usage

```sh
$ go vet -vettool=`which gencon` pkgname
```

## fixgencon

`fixgencon` is under development and supports only single type parameter.

```go
--- before `fixgencon` ---

package a

func f[T any](t T) {} // want "should not use 'any'. hint: string|~int"

func g[T, U any, V int](t T, u U, v V) {} // want "should not use 'any'. hint: bool|int" "should not use 'any'. hint: string|~int"

func invoke() { // OK
	f("gopher")
	f(1)
	type MyInt int
	f(MyInt(18))

	g(3, "gopher", 100)
	g(true, MyInt(3), 100)
}

--- after `fixgencon` ---

package a

func f[T string | ~int](t T) {} // want "should not use 'any'. hint: string|~int"

func g[T, U any, V int](t T, u U, v V) {} // want "should not use 'any'. hint: bool|int" "should not use 'any'. hint: string|~int"

func invoke() { // OK
	f("gopher")
	f(1)
	type MyInt int
	f(MyInt(18))

	g(3, "gopher", 100)
	g(true, MyInt(3), 100)
}
```

### install

```sh
$ go install github.com/ari1021/gencon/cmd/fixgencon@latest
```

### usage

```sh
$ fixgencon ./...
```
