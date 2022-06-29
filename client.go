package main

const clientBuffer = 5

type client struct {
	ID             int
	state          state
	lastInputFrame int
	inputBuffer    []input
}

func (c *client) run() {

	// send input

	// read frame

	//	if c.lastInputFrame < state.Frame+clientBuffer {
	// send input

	//	}

}
