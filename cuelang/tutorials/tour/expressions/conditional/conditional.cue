price: number

// Require a justification if price is too high
if price > 100 {
	justification: string
}

price: 200

a: {
	if price > 100 {
		test: "test"
	}
}
