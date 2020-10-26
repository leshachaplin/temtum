package data

type Status struct {
	Address string `db:"address"`
	Hash    string `db:"hash"`
	Type    string `db:"type"`
	Ban     bool   `db:"ban"`
	Stat    int    `db:"status"`
}
