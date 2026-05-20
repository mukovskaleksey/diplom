package datastruct

type Specialist struct {
	UserId      int64  `db:"user_id"`
	Category    string `db:"category"`
	CurrentLoad int    `db:"current_load"`
}
