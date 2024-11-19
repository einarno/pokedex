package pokedata

import (
	"math/rand"
)

type Pokemon struct {
	Name           string        `json:"name"`            // Pokémon's name
	ID             int           `json:"id"`              // Pokémon's ID
	Height         int           `json:"height"`          // Pokémon's height
	Weight         int           `json:"weight"`          // Pokémon's weight
	BaseExperience int           `json:"base_experience"` // Base experience gained when
	Types          []TypeElement `json:"types"`           // List of types
	Stats          []Stat        `json:"stats"`           // List of stats
}

type TypeElement struct {
	Slot int  `json:"slot"` // Order in which the type appears
	Type Info `json:"type"` // Details about the type
}

type Stat struct {
	BaseStat int  `json:"base_stat"` // Base value of the stat
	Effort   int  `json:"effort"`    // Effort points contributed by the stat
	Stat     Info `json:"stat"`      // Details about the stat
}

type Info struct {
	Name string `json:"name"` // Name of the entity (e.g., ability, type)
	URL  string `json:"url"`  // URL for more information
}

// getRandomBool generates a random boolean based on Base Experience
// Higher Base Experience increases the likelihood of False
func IsCatched(baseExperience int) bool {
	const maxBaseExperience = 255 // Highest Base Experience (Blissey)

	// Normalize base experience to a range [0, 1]
	normalized := float64(baseExperience) / float64(maxBaseExperience)


	randomValue := rand.Float64()

	// Return false if randomValue <= normalized, else true
	return randomValue > normalized
}
