package main

import (
	"fmt"
	"math/rand"
	"time"
)

/**Declare the Player "Class".*/
type Player struct {
	ID      string
	balance float64
}

/**Declare the bet "Class".*/
type Bet struct {
	PlayerID string
	Amount   float64
	Odds     float64
}

/**This function adds a specified amount to the player's balance, thus increasing the player's money limit.*/
func (p *Player) GetBalance() float64 {
	return p.balance
}

func (p *Player) SetBalance(newBalance float64) {
	p.balance = newBalance
}

/**This function deducts a specified amount from the player's balance if it is sufficient. If the player's balance is insufficient to pay, it will return an error message.*/
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

/**This function subtracts the specified bet amount from the player's balance. If the player's balance is not enough for the bet, it will return an error message.*/
func (p *Player) Deposit(amount float64) {
	p.SetBalance(p.GetBalance() + amount)
}

/**This function calculates the player's payout depending on the outcome of a bet.*/
func (p *Player) CalculatePayout(bet Bet, win bool) (float64, error) {
	if win {
		payout := bet.Amount * bet.Odds
		p.SetBalance(p.GetBalance() + payout)
		return payout, nil
	}
	return 0, nil
}

/** This function can be Calculation of multiplier based on the bet.*/
func CalculateOdds(betAmount float64) float64 {
	if betAmount <= 10 {
		return 1.2
	} else if betAmount <= 50 {
		return 1.5
	} else {
		return 2.0
	}
}

/** This function can be RTP calculation based on bet multiplier and winning chance.*/
func CalculateRTP(odds, winChance float64) float64 {
	return (winChance * odds) * 100
}

/**This Function can be create a new bet.*/
func NewBet(playerID string, amount float64) Bet {
	odds := CalculateOdds(amount)
	return Bet{
		PlayerID: playerID,
		Amount:   amount,
		Odds:     odds,
	}
}

/** This function is the process and business logic of the specific game.*/
func PlayGame(player *Player, bet Bet) (float64, error) {
	winChance := 0.495 / bet.Odds
	win := rand.Float64() < winChance

	payout, err := player.CalculatePayout(bet, win)
	if err != nil {
		return 0, err
	}
	return payout, nil
}

// Ask the user if they want to continue playing, validate input
func AskToContinue() bool {
	var response string
	for {
		fmt.Print("Do you want to continue playing? (yes/no): ")
		fmt.Scan(&response)
		if response == "yes" || response == "no" {
			break
		} else {
			fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
		}
	}
	return response == "yes"
}

// Ask the user to deposit more money, validate input
func AskToDeposit(player *Player) {
	var response string
	for {
		fmt.Print("You don't have enough money. Do you want to deposit more money? (yes/no): ")
		fmt.Scan(&response)
		if response == "yes" || response == "no" {
			break
		} else {
			fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
		}
	}

	if response == "yes" {
		var depositAmount float64
		fmt.Print("Enter deposit amount: ")
		fmt.Scan(&depositAmount)
		player.Deposit(depositAmount)
		fmt.Printf("New balance: %.2f\n", player.GetBalance())
	}
}

// Ask the user to place a bet with a minimum and maximum limit
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

// Test the program.
func main() {
	rand.Seed(time.Now().UnixNano())

	player := Player{ID: "player1"}
	player.SetBalance(100.0) // Direct adjustment of the balance with a setter.

	fmt.Printf("Starting balance: %.2f\n", player.GetBalance())

	minBet := 5.0 // Fix minimum tét

	for {
		fmt.Println("\n--- New Round ---")

		// We dynamically determine the maximum bet based on the player's current balance.
		maxBet := player.GetBalance()

		//If the balance is less than the minimum bet, we offer the deposit option.
		if maxBet < minBet {
			fmt.Println("Your balance is less than the minimum bet.")
			AskToDeposit(&player)

			maxBet = player.GetBalance()

			// If there is still not enough money, we end the game.
			if maxBet < minBet {
				fmt.Println("Still insufficient balance to place the minimum bet. Game over!")
				break
			}
		}

		// We ask for the bet from the player
		betAmount := AskForBetAmount(minBet, maxBet)
		bet := NewBet(player.ID, betAmount)
		fmt.Printf("Bet amount: %.2f, Odds: %.2f\n", bet.Amount, bet.Odds)

		// Bet processing
		err := player.PlaceBet(betAmount)
		if err != nil {
			fmt.Println("Error placing bet:", err)
			AskToDeposit(&player)
			continue
		}

		// Game round processing
		gamePayout, err := PlayGame(&player, bet)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Payout: %.2f\n", gamePayout)
		}

		fmt.Printf("Final balance: %.2f\n", player.GetBalance())

		if !AskToContinue() {
			fmt.Println("Thank you for playing!")
			break
		}
	}
}
