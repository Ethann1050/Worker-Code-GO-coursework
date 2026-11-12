for i := 0; i < p.Turns; i++ {
	worldSlice = calculateNextState(p, worldSlice)
	completedTurns = i + 1

	aliveCells = calculateAliveCells(p, worldSlice)
	c.events <- TurnComplete{CompletedTurns: completedTurns}

}