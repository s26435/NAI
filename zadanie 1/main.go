// Authors:
// Jan Wolski s26435
// Marcin Topolniak s25672

// Rules of game: https://pl.wikipedia.org/wiki/Nim

//Instrukcja przygotowania środowiska znajduje się w README.md w repozytorium

package main

import (
	"container/heap"
	"fmt"
	_ "math"
)

// Game represents a Nim game with multiple piles of stones.
type Game struct {
	piles []int
}

// A* Priority Queue Element
type Node struct {
	state    []int // The configuration of the piles
	cost     int   // g(n): Cost to reach this node (depth)
	priority int   // f(n): g(n) + h(n), priority in A*
	index    int   // The index of the node in the heap
}

// A Priority Queue implementation using heap.Interface
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// displayPiles prints the current state of the piles to the console.
func (g *Game) displayPiles() {
	fmt.Println("Current piles:")
	for i, pile := range g.piles {
		fmt.Printf("Pile %d: %d\n", i+1, pile)
	}
}

// isGameOver checks if the game is over, which occurs when all piles are empty.
func (g *Game) isGameOver() bool {
	for _, pile := range g.piles {
		if pile > 0 {
			return false
		}
	}
	return true
}

// removeStones updates the specified pile by removing the given number of stones.
func (g *Game) removeStones(pileIndex, stones int) {
	if pileIndex < 0 || pileIndex >= len(g.piles) || stones <= 0 || stones > g.piles[pileIndex] {
		fmt.Println("Invalid move!")
		return
	}
	g.piles[pileIndex] -= stones
}

// heuristic estimates the remaining cost (h(n)) to a winning position for AI.
// A simple heuristic is to count the total number of stones.
func (g *Game) heuristic(state []int) int {
	stoneCount := 0
	for _, pile := range state {
		stoneCount += pile
	}
	return stoneCount
}

// generateSuccessors returns all possible states (after AI move) from the current state.
func (g *Game) generateSuccessors(state []int) [][]int {
	var successors [][]int
	for i, pile := range state {
		for stones := 1; stones <= pile; stones++ {
			newState := make([]int, len(state))
			copy(newState, state)
			newState[i] -= stones
			successors = append(successors, newState)
		}
	}
	return successors
}

func (g *Game) aiMove() {
	fmt.Println("AI is making its move using A*...")
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	initialState := g.piles
	initialNode := &Node{
		state:    initialState,
		cost:     0,
		priority: g.heuristic(initialState),
	}
	heap.Push(&pq, initialNode)

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		currentNode := heap.Pop(&pq).(*Node)

		// Check if this state is game over (AI should leave all piles empty)
		if g.heuristic(currentNode.state) == 0 { // All piles are empty
			g.piles = currentNode.state
			fmt.Println("AI finished its move.")
			return
		}

		// Generate possible next moves (successor states)
		successors := g.generateSuccessors(currentNode.state)
		for _, successor := range successors {
			stateKey := fmt.Sprint(successor)
			if !visited[stateKey] {
				visited[stateKey] = true
				gScore := currentNode.cost + 1
				hScore := g.heuristic(successor)
				fScore := gScore + hScore

				// Push successor to priority queue
				newNode := &Node{
					state:    successor,
					cost:     gScore,
					priority: fScore,
				}
				heap.Push(&pq, newNode)

				// Set the game's piles to this state (AI's move)
				g.piles = successor
				fmt.Println("AI makes a move:")
				//g.displayPiles()
				return
			}
		}
	}
}

// playerMove is unchanged, allowing the player to make a move
func (g *Game) playerMove() {
	var pileIndex, stones int
	for {
		g.displayPiles()
		fmt.Print("Enter pile number (1, 2, ...): ")
		fmt.Scan(&pileIndex)
		fmt.Print("Enter number of stones to remove: ")
		fmt.Scan(&stones)
		if pileIndex > 0 && pileIndex <= len(g.piles) {
			if stones > 0 && stones <= g.piles[pileIndex-1] {
				break
			}
		}
		fmt.Println("Invalid move! Please try again.")
	}
	g.removeStones(pileIndex-1, stones)
}

// play starts the game loop, alternating between player and AI moves until the game is over.
func (g *Game) play() {
	for !g.isGameOver() {
		g.playerMove()
		if g.isGameOver() {
			fmt.Println("Game over! You win.")
			return
		}
		g.aiMove()
		if g.isGameOver() {
			fmt.Println("Game over! AI wins.")
			return
		}
	}
}

// main initializes the game with a starting configuration and begins play.
func main() {
	initialPiles := []int{3, 4, 5} // Example starting configuration
	game := Game{piles: initialPiles}
	game.play()
}
