a :: {
	foo?: int
	bar?: string
	baz?: string
}
b: a & {
	foo:  3
	baz?: 2 // baz?: _|_
}
