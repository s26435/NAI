package main
 
import (
	"fmt"
	"math/rand"
)
 
// Game represents a Nim game with multiple piles of stones.
type Game struct {
	piles []int // piles holds the number of stones in each pile.
}
 
// displayPiles prints the current state of the piles to the console.
func (g *Game) displayPiles() {
	fmt.Println("Current piles:")
	for i, pile := range g.piles {
		fmt.Printf("Pile %d: %d\n", i+1, pile)
	}
}
 
// isGameOver checks if the game is over, which occurs when all piles are empty.
// It returns true if the game is over, otherwise false.
func (g *Game) isGameOver() bool {
	for _, pile := range g.piles {
		if pile > 0 {
			return false
		}
	}
	return true
}
 
// removeStones updates the specified pile by removing the given number of stones.
// If the move is invalid (out of bounds or trying to remove more stones than available), it prints an error message.
func (g *Game) removeStones(pileIndex, stones int) {
	if pileIndex < 0 || pileIndex >= len(g.piles) || stones <= 0 || stones > g.piles[pileIndex] {
		fmt.Println("Invalid move!")
		return
	}
	g.piles[pileIndex] -= stones
}
 
// playerMove allows the player to make a move by selecting a pile and the number of stones to remove.
// It prompts the player until a valid move is made.
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
 
// aiMove allows the AI to make a move using a simple strategy.
// It tries to create a losing position for the player, or makes a random valid move if no optimal move is found.
func (g *Game) aiMove() {
	fmt.Println("AI is making its move...")
	for i := range g.piles {
		if g.piles[i] > 0 {
			// Basic strategy: remove stones to create a losing position for the player
			for stones := 1; stones <= g.piles[i]; stones++ {
				newPile := g.piles[i] - stones
				if g.calculateNimSum(newPile) == 0 {
					g.removeStones(i, stones)
					return
				}
			}
		}
	}
	// If no optimal move found, just make a random valid move
	for i := range g.piles {
		if g.piles[i] > 0 {
			stones := rand.Intn(g.piles[i]) + 1
			g.removeStones(i, stones)
			return
		}
	}
}
 
// calculateNimSum computes the Nim sum of the current game state.
// It returns the result of XORing all piles together.
func (g *Game) calculateNimSum(pile int) int {
	nimSum := pile
	for _, p := range g.piles {
		nimSum ^= p
	}
	return nimSum
}
 
// play starts the game loop, alternating between player and AI moves until the game is over.
func (g *Game) play() {
	for !g.isGameOver() {
		g.playerMove()
		if g.isGameOver() {
			fmt.Println("Game over! You Win.")
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