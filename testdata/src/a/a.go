package a // want "error"

func f[T any](t T) {}

func invoke() {
	f("gopher")
	f(1)
	type MyInt int
	f(MyInt(18))
}
