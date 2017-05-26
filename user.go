package main

import (
	"net/http"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

type User struct {
	gorm.Model
	Email  string `json:"email" gorm:"unique_index;"`
	Tokens []Token
}

type Token struct {
	gorm.Model
	UserID       int `gorm:"index"`
	Service      string
	AccessToken  string
	ExpiryToken  string
	RefreshToken string
	TokenType    string
}

type UserORM struct {
	db *gorm.DB
}

func userRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", createUserRequest)
	// r.Get("/user/:userID", u.getUser)
	return r
}

func createUserRequest(w http.ResponseWriter, r *http.Request) {
	userR := &UserRequest{}
	if err := render.Bind(r, userR); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user := userR.User
	db, _ := DB()
	// defer db.Close()
	if db.Find(&user, "email = ?", user.Email).RecordNotFound() == false {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("user already exists")))
		db.Close()
		return
	}
	if err := db.Create(&user).Error; err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		db.Close()
		return
	}
	db.Close()
	fmt.Println(user.Email)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

type UserRequest struct {
	*User
}

func (u *UserRequest) Bind(r *http.Request) error {
	return nil
}

type UserResponse struct {
	Email string `json:"email"`
}

func NewUserResponse(user *User) *UserResponse {
	resp := &UserResponse{Email: user.Email}
	return resp
}

func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
