package auth

import (
	"bytes"
	"database/sql"
	"errors"
	"golang.org/x/crypto/argon2"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	ID      uint32
	Login   string
	IsAdmin bool
}

type UserHandler struct {
	DB       *sql.DB
	Tmpl     *template.Template
	Sessions SessionManager
}

// @Summary Аутентифицирует пользователя
// @Description Аутентифицирует пользователя по логину и паролю, предоставленным в запросе.
// @Accept json
// @Produce json
// @Param login query string true "Логин пользователя"
// @Param password query string true "Пароль пользователя"
// @Success 200 {string} string "Logged in"
// @Failure 400 {string} string "No user" "Bad pass"
// @Failure 500 {string} string "Db err"
// @Router /login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("password")

	user, err := uh.checkPasswordByLogin(login, pass)

	switch err {
	case nil:
		// all is ok
	case errNoRec:
		http.Error(w, "No user", http.StatusBadRequest)
	case errBadPass:
		http.Error(w, "Bad pass", http.StatusBadRequest)
	default:
		http.Error(w, "Db err", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}
	uh.Sessions.Create(w, user)
	w.Write([]byte("Logged in"))
	w.WriteHeader(http.StatusOK)
	//http.Redirect(w, r, "/user", http.StatusFound)
}

// @Summary Регистрирует нового пользователя
// @Description Регистрирует нового пользователя, предоставляя логин, пароль и статус isAdmin.
// @Accept json
// @Produce json
// @Param login query string true "Логин пользователя"
// @Param password query string true "Пароль пользователя"
// @Param isAdmin query boolean true "Статус isAdmin пользователя (true или false)"
// @Success 201 {string} string "User created"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	salt := RandStringRunes(8)
	pass := uh.hashPass(r.FormValue("password"), salt)
	isAdminString := r.FormValue("isAdmin")
	var isAdmin bool
	if isAdminString == "true" {
		isAdmin = true
	} else if isAdminString == "false" {
		isAdmin = false
	} else {
		http.Error(w, "isAdmin must be true or false", http.StatusBadRequest)
		return
	}
	result, err := uh.DB.Exec("INSERT INTO users(login, password, role) VALUES($1, $2, $3)", login, pass, isAdmin)
	if err != nil {
		log.Println("insert error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
		return
	}
	userID, _ := result.LastInsertId()

	user := &User{
		ID:      uint32(userID),
		IsAdmin: isAdmin,
		Login:   login,
	}
	uh.Sessions.Create(w, user)
	w.Write([]byte("User created"))
	w.WriteHeader(http.StatusCreated)
	//http.Redirect(w, r, "/user", http.StatusFound)
}

// @Summary Выход пользователя
// @Description Разлогинивает текущего пользователя, уничтожая его сессию, и перенаправляет на страницу входа.
// @Accept json
// @Produce json
// @Success 200 {string} string "Logged out"
// @Router /logout [get]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	uh.Sessions.DestroyCurrent(w, r)
	w.Write([]byte("Logged out"))
	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (uh *UserHandler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt[:len(salt)])
	return append(res, hashedPass...)
}

var (
	errNoRec   = errors.New("No user record found")
	errBadPass = errors.New("No user record found")
)

func (uh *UserHandler) passwordIsValid(pass string, row *sql.Row) (*User, error) {

	var (
		dbPass []byte
		user   = &User{}
	)
	err := row.Scan(&user.ID, &user.Login, &user.IsAdmin, &dbPass)
	if err == sql.ErrNoRows {
		return nil, errNoRec
	} else if err != nil {
		return nil, err
	}

	salt := string(dbPass[0:8])
	if !bytes.Equal(uh.hashPass(pass, salt), dbPass) {
		return nil, errBadPass
	}
	return user, nil
}

func (uh *UserHandler) checkPasswordByUserID(uid uint32, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, role, password FROM users WHERE id = $1", uid)
	return uh.passwordIsValid(pass, row)
}

func (uh *UserHandler) checkPasswordByLogin(login, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, role, password FROM users WHERE login = $1", login)
	return uh.passwordIsValid(pass, row)
}

func (uh *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pass1") == "" || r.FormValue("pass1") != r.FormValue("pass2") {
		http.Error(w, "New password mistmatch", http.StatusBadRequest)
		return
	}

	sess, _ := SessionFromContext(r.Context())
	user, err := uh.checkPasswordByUserID(sess.UserID, r.FormValue("old_password"))
	if err != nil {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	salt := RandStringRunes(8)
	pass := uh.hashPass(r.FormValue("pass1"), salt)

	_, err = uh.DB.Exec("UPDATE users SET password = $1 WHERE id = $2",
		pass, user.ID)
	if err != nil {
		log.Println("update password error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	uh.Sessions.DestroyAll(w, user)
	uh.Sessions.Create(w, user)

	//http.Redirect(w, r, "/", http.StatusFound)
}
