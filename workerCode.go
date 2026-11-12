package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type CalculateTurnService struct{}

type Request struct {
	WorldSlice  [][]byte
	Turns       int
	ImageHeight int
	ImageWidth  int
}

type Response struct {
	NewWorldSlice  [][]byte
	CompletedTurns int
}

func calculateNextState(height, width int, world [][]byte) [][]byte {
	newWorld := make([][]byte, height)
	for i := range newWorld {
		newWorld[i] = make([]byte, width)
	}

	for row := 0; row < height; row++ {
		for column := 0; column < width; column++ {
			liveCells := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nRow := (row + dr + height) % height
					nCol := (column + dc + width) % width
					if world[nRow][nCol] == 255 {
						liveCells++
					}
				}
			}
			if world[row][column] == 255 {
				if liveCells == 2 || liveCells == 3 {
					newWorld[row][column] = 255
				}
			} else if liveCells == 3 {
				newWorld[row][column] = 255
			}
		}
	}
	return newWorld
}

func (ct *CalculateTurnService) CalculateTurns(request Request, response *Response) error {
	world := request.WorldSlice
	completedTurns := 0

	for i := 0; i < request.Turns; i++ {
		world = calculateNextState(request.ImageHeight, request.ImageWidth, world)
		completedTurns = i + 1
	}

	response.NewWorldSlice = world
	response.CompletedTurns = completedTurns

	fmt.Printf("Processed %d turns for %dx%d grid\n", completedTurns, request.ImageWidth, request.ImageHeight)
	return nil
}

func main() {
	service := new(CalculateTurnService)
	err := rpc.Register(service)
	if err != nil {
		fmt.Println("Failed to register RPC: ", err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Failed to listen: ", err)
	}
	defer listener.Close()

	fmt.Println("Listening on 0.0.0.0:8080...")
	rpc.Accept(listener)
}
