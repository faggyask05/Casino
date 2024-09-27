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

// Test the program.
func main() {

	rand.Seed(time.Now().UnixNano())

	player := Player{ID: "player1"}
	player.SetBalance(100.0) // A balance közvetlen beállítása setterrel

	fmt.Printf("Starting balance: %.2f\n", player.GetBalance())

	player.Deposit(50)
	fmt.Printf("After deposit: %.2f\n", player.GetBalance())

	// Create a new bet with an adjustable stake
	betAmount := 150.0
	bet := NewBet(player.ID, betAmount)
	fmt.Printf("Bet amount: %.2f, Odds: %.2f\n", bet.Amount, bet.Odds)

	//RTP Calculation.
	rtp := CalculateRTP(bet.Odds, 0.495)
	fmt.Println("RTP: ", rtp)

	//Testing with more rounds.
	for i := 0; i < 10; i++ {
		fmt.Printf("Round %d:\n", i+1)

		// Placing a bet
		err := player.PlaceBet(bet.Amount)
		if err != nil {
			fmt.Println("Error placing bet:", err)
			break // If not enought mony for the player, the game is ending.
		}

		gamePayout, err := PlayGame(&player, bet)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Payout: %.2f\n", gamePayout)
		}
		fmt.Printf("Final balance: %.2f\n", player.GetBalance())

		if player.GetBalance() < bet.Amount {
			fmt.Println("Insufficient balance to continue playing.") //We will check if there is still enough balance to continue
			break
		}
	}
}
