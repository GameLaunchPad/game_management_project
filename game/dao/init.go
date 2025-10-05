package dao

type gameDAO struct{}

// NewGameDAO creates a new GameDAO.
func NewGameDAO() IGameDAO {
	return &gameDAO{}
}
