package data

import "strconv"

func generateXpLeaderboardKey(limit int) string {
	return "xp_leaderboard_" + strconv.FormatInt(int64(limit), 10)
}
