package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var sessions = struct {
	sync.RWMutex
	m map[string]int // token -> userID
}{m: make(map[string]int)}

func hashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password too long")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func createSession(userID int) (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}
	sessions.Lock()
	sessions.m[token] = userID
	sessions.Unlock()
	return token, nil
}

func getUserIDByToken(token string) (int, bool) {
	sessions.RLock()
	defer sessions.RUnlock()
	id, ok := sessions.m[token]
	return id, ok
}

func deleteSession(token string) {
	sessions.Lock()
	delete(sessions.m, token)
	sessions.Unlock()
}
