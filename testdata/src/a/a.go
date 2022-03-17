package a

func e[T any](t T) {} // want "change constraint of T from any to string|~int"

func f[T, U any, V int](t T, u U, v V) {} // want "change constraint of T from any to bool|int" "change constraint of U from any to string|~int"

func g(s string) {} // OK

func invoke() { // OK
	e("gopher")
	e(1)
	type MyInt int
	e(MyInt(18))

	f(3, "gopher", 100)
	f(true, MyInt(3), 100)

	g("gopher")
}
