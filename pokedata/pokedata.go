package pokedata

import "math/rand/v2"

type Pokemon struct {
	ID             int    `json:"id"`              // Pokémon ID
	Name           string `json:"name"`            // Pokémon name
	BaseExperience *int   `json:"base_experience"` // Base experience (nullable)
	Height         int    `json:"height"`          // Height (nullable)
	Weight         int    `json:"weight"`          // Weight (nullable)
	Types          []Type `json:"types"`           // Pokémon types
	Stats          []Stat `json:"stats"`           // Pokémon stats
}

type Type struct {
	Slot int  `json:"slot"` // Slot for the type
	Type Info `json:"type"` // Type details (name and URL)
}

type Stat struct {
	BaseStat int  `json:"base_stat"` // Base value of the stat
	Stat     Info `json:"stat"`      // Details about the stat
}

type Info struct {
	Name string `json:"name"` // Name of the entity
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
