package Utils

import (
	"fmt"
	"sort"
	"xxx/shared"
)

func SortUserScoresByScoreDesc(scores []shared.UserScore) []shared.UserScore {
	fmt.Println(scores)
	oldPlaces := make(map[string]int)
	for _, user := range scores {
		oldPlaces[user.UserId] = user.Place
	}

	sorted := make([]shared.UserScore, len(scores))
	copy(sorted, scores)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].TotalScore > sorted[j].TotalScore
	})

	for i := range sorted {
		user := &sorted[i]
		newPlace := i + 1
		oldPlace := oldPlaces[user.UserId]
		fmt.Printf("%v,%v,%v\n", user.UserId, oldPlace, newPlace)
		if oldPlace != 0 {
			if oldPlace > newPlace {
				user.Progress = true
			} else {
				user.Progress = false
			}
		} else {
			user.Progress = false
		}
		user.Place = newPlace
	}
	return sorted
}
