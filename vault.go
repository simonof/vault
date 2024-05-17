package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Secret struct {
	Value     string
	TTL       time.Duration
	CreatedAt time.Time
}

type User struct {
	Username string
	Password string
}

type AppRole struct {
	RoleID   string
	SecretID string
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type Vault struct {
	Secrets  map[string]Secret
	Users    map[string]User
	AppRoles map[string]AppRole
	Tokens   map[string]Token
	Mutex    sync.Mutex
}

func NewVault() *Vault {
	return &Vault{
		Secrets:  make(map[string]Secret),
		Users:    make(map[string]User),
		AppRoles: make(map[string]AppRole),
		Tokens:   make(map[string]Token),
	}
}

func (v *Vault) CreateRoot(name string, secretValue string, ttl time.Duration) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	v.Secrets[name] = Secret{Value: secretValue, TTL: ttl, CreatedAt: time.Now()}
}

func (v *Vault) CreateUser(username string, password string) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	v.Users[username] = User{Username: username, Password: password}
}

func (v *Vault) CreateAppRole(roleID string, secretID string) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	v.AppRoles[roleID] = AppRole{RoleID: roleID, SecretID: secretID}
}

func (v *Vault) GenerateToken() (string, error) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	tokenValue := generateRandomString(32)
	v.Tokens[tokenValue] = Token{Value: tokenValue, ExpiresAt: time.Now().Add(24 * time.Hour)}
	return tokenValue, nil
}

func (v *Vault) GetTokenByAppRole(roleID string, secretID string) (string, error) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	appRole, exists := v.AppRoles[roleID]
	if !exists || appRole.SecretID != secretID {
		return "", errors.New("invalid roleID or secretID")
	}
	return v.GenerateToken()
}

func (v *Vault) GetSecretByToken(token string, rootName string) (string, error) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	tok, exists := v.Tokens[token]
	if !exists || time.Now().After(tok.ExpiresAt) {
		return "", errors.New("invalid or expired token")
	}
	secret, exists := v.Secrets[rootName]
	if !exists || time.Now().After(secret.CreatedAt.Add(secret.TTL)) {
		return "", errors.New("secret not found or expired")
	}
	return secret.Value, nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	vault := NewVault()
	http.HandleFunc("/create-root", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")
		ttl, _ := time.ParseDuration(r.URL.Query().Get("ttl"))
		vault.CreateRoot(name, value, ttl)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Root created"))
	})
	http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		vault.CreateUser(username, password)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User created"))
	})
	http.HandleFunc("/create-approle", func(w http.ResponseWriter, r *http.Request) {
		roleID := r.URL.Query().Get("roleID")
		secretID := r.URL.Query().Get("secretID")
		vault.CreateAppRole(roleID, secretID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AppRole created"))
	})
	http.HandleFunc("/get-token", func(w http.ResponseWriter, r *http.Request) {
		roleID := r.URL.Query().Get("roleID")
		secretID := r.URL.Query().Get("secretID")
		token, err := vault.GetTokenByAppRole(roleID, secretID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	})
	http.HandleFunc("/get-secret", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		rootName := r.URL.Query().Get("rootName")
		secret, err := vault.GetSecretByToken(token, rootName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(secret))
	})

	fmt.Println("Vault is running on :8080")
	http.ListenAndServe(":8080", nil)
}
