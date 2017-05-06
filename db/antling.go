package db

type Antling struct {
	Id int `json:"id"`
}

const antling_save_query = `
INSERT INTO anthive.antling
DEFAULT VALUES
RETURNING anthive.antling.id
`

func (a *Antling) Save() error {
	return conn.QueryRow(antling_save_query).Scan(&a.Id)
}
