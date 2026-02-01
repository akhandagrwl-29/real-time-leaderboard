// real-time-leaderboard/leaderboard.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	ctx            = context.Background()
	redisClient    *redis.Client
	leaderboardKey = "global_leaderboard_%s"
)

type ScoreSubmission struct {
	User  string  `json:"user"`
	Score float64 `json:"score"`
	Game  string  `json:"game"`
}

func submitScore(w http.ResponseWriter, r *http.Request) {
	var submission ScoreSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// Add or update user's score in the sorted set
	err := redisClient.ZAdd(ctx, fmt.Sprintf(leaderboardKey, submission.Game), redis.Z{
		Score:  submission.Score,
		Member: submission.User,
	}).Err()
	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Score submitted for user %s\n", submission.User)
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	topN := 10
	if n := r.URL.Query().Get("n"); n != "" {
		if val, err := strconv.Atoi(n); err == nil {
			topN = val
		}
	}
	game := r.URL.Query().Get("game")
	leaderboardKey := fmt.Sprintf(leaderboardKey, game)

	results, err := redisClient.ZRevRangeWithScores(ctx, leaderboardKey, 0, int64(topN-1)).Result()
	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	leaderboard := make([]map[string]interface{}, 0, len(results))
	for i, z := range results {
		leaderboard = append(leaderboard, map[string]interface{}{
			"rank":  i + 1,
			"user":  z.Member,
			"score": z.Score,
			"game":  game,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

func getUserRank(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "Missing user parameter", http.StatusBadRequest)
		return
	}
	game := r.URL.Query().Get("game")
	leaderboardKey := fmt.Sprintf(leaderboardKey, game)

	rank, err := redisClient.ZRevRank(ctx, leaderboardKey, user).Result()
	if err == redis.Nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	score, _ := redisClient.ZScore(ctx, leaderboardKey, user).Result()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"rank":  rank + 1,
		"score": score,
		"game":  game,
	})
}

func uiMain(w http.ResponseWriter, r *http.Request) {
	mainUITmpl.Execute(w, nil)
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	http.HandleFunc("/submit", submitScore)
	http.HandleFunc("/leaderboard", getLeaderboard)
	http.HandleFunc("/rank", getUserRank)

	http.HandleFunc("/ui", uiMain)

	fmt.Println("Leaderboard server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
