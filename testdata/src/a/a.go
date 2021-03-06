package a

func e[T any](t T) {} // want "should not use 'any'. hint: string|~int"

func f[T, U any, V int](t T, u U, v V) {} // want "should not use 'any'. hint: bool|int" "should not use 'any'. hint: string|~int"

func g(s string) {} // OK

func h[T any](t T) {} // want "should not use 'any'"

func invoke() { // OK
	e("gopher")
	e(1)
	type MyInt int
	e(MyInt(18))

	f(3, "gopher", 100)
	f(true, MyInt(3), 100)

	g("gopher")
}
