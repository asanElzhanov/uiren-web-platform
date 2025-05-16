package progress

/*
updateUserProgressRequest json

	{
		"user_id": "string",
		"xp": 0,
		"new_badges":[
		"badge1",
		"badge2"
		],
		"achievements_progress":[
		{
			"achievement_id": 1,
			"earned_progress": 52
		},
		{
			"achievement_id": 2,
			"earned_progress": 100
		}
			]
	}
*/

// for update
type UserAchievement struct {
	AchievementName  string `json:"ach_name"`
	Level            int    `json:"lvl"`
	LevelDescription string `json:"description"`
	Progress         int    `json:"progress"`
	Threshold        int    `json:"lvl_threshold"`
}

type AchievementProgress struct {
	AchievementID  int `json:"achievement_id"`
	EarnedProgress int `json:"earned_progress"`
	NewLevel       int
}

type UpdateUserProgressRequest struct {
	UserID               string                `json:"user_id"`
	XP                   int                   `json:"xp"`
	NewBadges            []string              `json:"new_badges"`
	AchievementsProgress []AchievementProgress `json:"achievements_progress"`
}

type AddBadgesRequest struct {
	UserID string   `json:"user_id"`
	Badges []string `json:"badges"`
}

type AddXPRequest struct {
	UserID string `json:"user_id"`
	XP     int    `json:"xp"`
}

type UpdateAchievementProgressRequest struct {
	UserID   string              `json:"user_id"`
	Progress AchievementProgress `json:"achievements_progress"`
}

// for insert badges
type Badge struct {
	Badge       string `json:"badge"`
	Description string `json:"description"`
}

// for leaderboard

type XPLeaderboard struct {
	Leaders []XPLeaderboardEntry `json:"leaders"`
	Total   int                  `json:"total"`
}

type XPLeaderboardEntry struct {
	Rank     int    `json:"rank"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	XP       int    `json:"xp"`
}
