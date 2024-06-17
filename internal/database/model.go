package database

type Task struct {
	Id         int64  `json:"id"`
	Belongs_to int64  `json:"belongs_to"`
	Order      int64  `json:"order"`
	Title      string `json:"title"`
	Details    string `json:"details"`
	State      int64  `json:"state"`
	Due        string `json:"due"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username" binding:"min=2,max=20,required"`
	Password string `json:"password" binding:"min=2,required"`
}

type Categories struct {
	Id         int64  `json:"id"`
	Belongs_to int64  `json:"belongs_to"`
	Name       string `json:"name"`
	Order      int64  `json:"order"`
}
