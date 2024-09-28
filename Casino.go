package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Define the Player struct.
type Player struct {
	ID      string
	balance float64
}

// Define the Bet struct.
type Bet struct {
	PlayerID string
	Amount   float64
	Odds     float64
}

// GetBalance returns the player's current balance.
func (p *Player) GetBalance() float64 {
	return p.balance
}

// SetBalance sets a new balance for the player.
func (p *Player) SetBalance(newBalance float64) {
	p.balance = newBalance
}

// PlaceBet deducts the bet amount from the player's balance.
func (p *Player) PlaceBet(betAmount float64) error {
	if betAmount <= 0 {
		return fmt.Errorf("invalid bet amount")
	}
	if p.GetBalance() < betAmount {
		return fmt.Errorf("insufficient balance")
	}
	p.SetBalance(p.GetBalance() - betAmount)
	return nil
}

// Deposit adds money to the player's balance.
func (p *Player) Deposit(amount float64) {
	p.SetBalance(p.GetBalance() + amount)
}

// CalculatePayout calculates the player's payout based on the bet and whether they win.
func (p *Player) CalculatePayout(bet Bet, win bool) (float64, error) {
	if win {
		payout := bet.Amount * bet.Odds
		p.SetBalance(p.GetBalance() + payout)
		return payout, nil
	}
	return 0, nil
}

// CalculateOdds computes the odds based on the RTP constant and bet amount.
// This is the first-degree equation where payoutMultiplier = RTP * prediction.
func CalculateOdds(rtp, prediction float64) float64 {
	return rtp * prediction
}

// GenerateCryptoRandom generates a secure random float between 0 and 1.
func GenerateCryptoRandom() (float64, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0, err
	}
	return float64(nBig.Int64()) / 1000000.0, nil
}

// NewBet creates a new bet based on the player's choice and bet amount.
func NewBet(playerID string, amount float64, rtp float64) Bet {
	odds := CalculateOdds(rtp, amount)
	return Bet{
		PlayerID: playerID,
		Amount:   amount,
		Odds:     odds,
	}
}

// PlayGame simulates a round of the game.
func PlayGame(player *Player, bet Bet) (float64, error) {
	// Calculate winChance based on bet amount
	winChance := 1 / (1 + bet.Amount/10)

	// Secure random generator using crypto/rand
	winRandom, err := GenerateCryptoRandom()
	if err != nil {
		return 0, err
	}

	// Determine if player wins based on winChance
	win := winRandom < winChance

	// Calculate payout based on whether the player won or lost
	payout, err := player.CalculatePayout(bet, win)
	if err != nil {
		return 0, err
	}

	return payout, nil
}

// AskToContinue prompts the player if they want to continue playing.
func AskToContinue() bool {
	var response string
	for {
		fmt.Print("Do you want to continue playing? (y/n): ")
		fmt.Scan(&response)
		if response == "y" || response == "n" {
			break
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
	return response == "y"
}

// AskToDeposit prompts the player to deposit more money.
func AskToDeposit(player *Player) {
	var response string
	for {
		fmt.Print("You don't have enough money. Do you want to deposit more money? (y/n): ")
		fmt.Scan(&response)
		if response == "y" || response == "n" {
			break
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}

	if response == "y" {
		var depositAmount float64
		fmt.Print("Enter deposit amount: ")
		fmt.Scan(&depositAmount)
		player.Deposit(depositAmount)
		fmt.Printf("New balance: %.2f\n", player.GetBalance())
	}
}

// AskForBetAmount prompts the player to place a bet with a min and max limit.
func AskForBetAmount(minBet, maxBet float64) float64 {
	var betAmount float64
	for {
		fmt.Printf("Enter your bet amount (minimum: %.2f, maximum: %.2f): ", minBet, maxBet)
		fmt.Scan(&betAmount)
		if betAmount >= minBet && betAmount <= maxBet {
			break
		} else {
			fmt.Println("Invalid bet amount. Please enter a value within the allowed range.")
		}
	}
	return betAmount
}

// CalculateAverageRTP simulates 1 million iterations to find the average RTP output.
func CalculateAverageRTP(rtp float64, iterations int) float64 {
	totalPayout := 0.0

	for i := 0; i < iterations; i++ {
		// Simulate a random bet amount
		betAmount := float64(5 + (i % 10)) // Example bet amount
		odds := CalculateOdds(rtp, betAmount)

		// Simulate win or lose (50% chance)
		winRandom, _ := GenerateCryptoRandom()
		win := winRandom < 0.5 // 50% win chance

		if win {
			totalPayout += betAmount * odds
		} else {
			totalPayout -= betAmount
		}
	}

	return totalPayout / float64(iterations)
}

// Main function to simulate the game.
func main() {
	rtpConstant := 0.95 // Example RTP constant
	player := Player{ID: "player1"}
	player.SetBalance(100.0) // Start with a balance of 100

	fmt.Printf("Starting balance: %.2f\n", player.GetBalance())

	minBet := 5.0 // Set minimum bet

	for {
		fmt.Println("\n--- New Round ---")

		maxBet := player.GetBalance()

		if maxBet < minBet {
			fmt.Println("Your balance is less than the minimum bet.")
			AskToDeposit(&player)

			maxBet = player.GetBalance()
		}

		if maxBet < minBet {
			fmt.Println("Still insufficient balance to place the minimum bet. Game over!")
			break
		}

		betAmount := AskForBetAmount(minBet, maxBet)
		bet := NewBet(player.ID, betAmount, rtpConstant)
		fmt.Printf("Bet amount: %.2f, Odds: %.2f\n", bet.Amount, bet.Odds)

		err := player.PlaceBet(betAmount)
		if err != nil {
			fmt.Println("Error placing bet:", err)
			AskToDeposit(&player)
			continue
		}

		gamePayout, err := PlayGame(&player, bet)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Payout: %.2f\n", gamePayout)
		}

		fmt.Printf("Final balance: %.2f\n", player.GetBalance())

		if !AskToContinue() {
			fmt.Printf("Your final balance: %.2f\n", player.GetBalance())
			fmt.Println("Thank you for playing!")
			break
		}
	}
	// Calculate average RTP output
	averageRTP := CalculateAverageRTP(rtpConstant, 1000000)
	fmt.Printf("Average RTP output over 1 million iterations: %.2f\n", averageRTP)

}
