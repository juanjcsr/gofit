package main

import (
	"context"
	"net/http"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique_index;"`
	Username string `json:"username" gorm:"unique_index"`
	Tokens   []Token
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

type userCtxKey string

var (
	contextUserKey = userCtxKey("token")
)

func (u userCtxKey) String() string {
	return "user ctx " + string(u)
}

func userRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/", createUserRequest)
	r.Route("/token/:username", func(r chi.Router) {
		r.Use(UserCtx)
		r.Get("/", getOrCreateUserToken)
	})

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

func getOrCreateUserToken(w http.ResponseWriter, r *http.Request) {
	// user := r.Context().Value(contextUserKey).(uint)
	u := r.Context().Value(contextUserKey).(*User)
	fmt.Printf("UUUUUUUU %v", u)
	var tokens []Token
	var count int
	db, _ := DB()
	defer db.Close()
	db.Model(&u).Related(&tokens).Count(count)
	if count == 0 {
		// show url to fitbit
	}
	fmt.Fprintf(w, "golaaa %v", tokens)
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userEmail := chi.URLParam(r, "username")
		fmt.Printf("username: %s\n", userEmail)
		var user User
		db, _ := DB()
		defer db.Close()
		if db.First(&user, "username = ?", userEmail).RecordNotFound() {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("no such user %s", userEmail)))
			return
		}
		fmt.Printf("usuario: %v", user)
		//ctx := context.WithValue(r.Context(), contextUserKey, user.ID)
		ctx := context.WithValue(r.Context(), contextUserKey, &user)
		// ctx = context.WithValue(ctx, "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type UserRequest struct {
	*User
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
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
