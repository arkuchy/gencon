package a

func f[T any](t T) {} // want "change any to string|~int"

func g[T, U any, V int](t T, u U, v V) {}

func invoke() {
	f("gopher")
	f(1)
	type MyInt int
	f(MyInt(18))

	g(3, "hoge", 100)
}
