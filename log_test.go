package log

func Example_Test() {
	SetLevel(DebugLevel)

	Debug("hello")
	Infow("hello")
	Warn("hello")
	//output:
}

func Example_Error() {
	SetLevel(DebugLevel)

	Error("hello")
	//output:
}

func Example_Panic() {
	SetLevel(DebugLevel)
	defer recover()

	Panic("hello")
	//output:
}

func Example_Fatal() {
	SetLevel(DebugLevel)

	Fatal("hello")
	//output:
}
