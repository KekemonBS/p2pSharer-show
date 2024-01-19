package tracker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KekemonBS/p2pSharer-back/tracker/models"
	"github.com/redis/go-redis/v9"
)

type CacheImpl interface {
	SaveUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, u models.User) error
	ReadUser(ctx context.Context, username string) (models.User, error)
	Save(ctx context.Context, u models.User, f models.Folder) error
	Read(ctx context.Context, u models.User) ([]models.Folder, error)
	Delete(ctx context.Context, u models.User, folderName string) error
}

type CacheHandlersImpl struct {
	ctx    context.Context
	logger *log.Logger
	ci     CacheImpl
}

func New(ctx context.Context, logger *log.Logger, ci CacheImpl) *CacheHandlersImpl {
	return &CacheHandlersImpl{
		ctx:    ctx,
		logger: logger,
		ci:     ci,
	}
}

// One table

func (c CacheHandlersImpl) SaveUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	username := r.URL.Query().Get("username")
	user.Name = username
	password := r.URL.Query().Get("password")
	providedHash := sha256.Sum256([]byte(password))
	user.PassHash = providedHash
	c.logger.Println(providedHash, " -- ", user.PassHash)
	if username == "" || password == "" {
		c.logger.Println("not enough/malformed parameters")
		http.Error(w, "not enough/malformed parameters", http.StatusBadRequest)
		return
	}

	// Check if user exists
	_, err := c.ci.ReadUser(c.ctx, user.Name)
	if err != nil && err != redis.Nil {
		c.logger.Println(err)
		http.Error(w, "User already exists", http.StatusInternalServerError)
		return
	}

	if err != nil {
		// Save new user
		err = c.ci.SaveUser(c.ctx, user)
		if err != nil {
			c.logger.Println(err)
			http.Error(w, "Error processing SaveUser", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		// User already exists
		w.WriteHeader(http.StatusConflict)
	}

	c.logger.Printf("Saved %s\n", user)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user) // Return the saved or existing user
}

func (c CacheHandlersImpl) ReadUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		c.logger.Println("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := c.ci.ReadUser(c.ctx, username)
	if err != nil {
		c.logger.Println(err)
		if err == redis.Nil {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error processing ReadUser", http.StatusInternalServerError)
		}
		return
	}

	c.logger.Printf("Read %s\n", username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c CacheHandlersImpl) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := []byte(r.URL.Query().Get("password"))
	providedHash := sha256.Sum256(password)

	c.logger.Println(password)
	if username == "" || password == nil {
		c.logger.Println("Username and hash are required")
		http.Error(w, "Username and hash are required", http.StatusBadRequest)
		return
	}

	// Retrieve user from Redis
	user, err := c.ci.ReadUser(c.ctx, username)
	c.logger.Println(user)
	if err != nil {
		c.logger.Println(err)
		if err == redis.Nil {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error processing DeleteUser", http.StatusInternalServerError)
		}
		return
	}

	// Validate password hash
	c.logger.Println("Got password :  ", password)
	c.logger.Println("Got password :  ", providedHash)
	if !checkHash(providedHash, user.PassHash) {
		c.logger.Println(providedHash, " -- ", user.PassHash)
		c.logger.Println("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Delete user from Redis
	err = c.ci.DeleteUser(c.ctx, models.User{Name: username, PassHash: providedHash})
	if err != nil {
		c.logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.logger.Printf("Deleted %s\n", user)
	w.WriteHeader(http.StatusOK)
}

func (c CacheHandlersImpl) Save(w http.ResponseWriter, r *http.Request) {
	var folder models.Folder
	err := json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		c.logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	password := []byte(r.URL.Query().Get("password"))
	providedHash := sha256.Sum256(password)

	// Validate user credentials
	user, err := c.ci.ReadUser(c.ctx, username)
	if err != nil {
		c.logger.Println(err)
		if err == redis.Nil {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error processing DeleteUser", http.StatusInternalServerError)
		}
		return
	}

	if !checkHash(providedHash, user.PassHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Save folder to Redis
	err = c.ci.Save(c.ctx, user, folder)
	if err != nil {
		c.logger.Println(err)
		http.Error(w, "Error saving folder", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	c.logger.Printf("Saved folder %s\n", folder.Name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(folder)
}

func (c CacheHandlersImpl) Read(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := []byte(r.URL.Query().Get("password"))
	providedHash := sha256.Sum256(password)

	if username == "" || password == nil {
		c.logger.Println("Not enough/malformed parameters")
		http.Error(w, "Not enough/malformed parameters", http.StatusBadRequest)
		return
	}

	// Check credentials
	user, err := c.ci.ReadUser(c.ctx, username)
	if err != nil {
		c.logger.Println(err)
		if err == redis.Nil {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error processing Read", http.StatusInternalServerError)
		}
		return
	}
	if !checkHash(providedHash, user.PassHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	folders, err := c.ci.Read(c.ctx, user)
	if err != nil {
		c.logger.Println(err)
		if err == redis.Nil {
			http.Error(w, "No folders found", http.StatusNotFound)
		} else {
			http.Error(w, "Error processing Read", http.StatusInternalServerError)
		}
		return
	}

	c.logger.Printf("Red folders for %s\n", user)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(folders)
}

func (c CacheHandlersImpl) Delete(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	folderName := r.URL.Query().Get("folderName")

	password := []byte(r.URL.Query().Get("password"))
	providedHash := sha256.Sum256(password)

	if username == "" || folderName == "" || password == nil {
		http.Error(w, "Username, folder name, and hash are required", http.StatusBadRequest)
		return
	}

	// Check credentials
	user, err := c.ci.ReadUser(c.ctx, username)
	if err != nil {
		c.logger.Println(err)
	}

	if !checkHash(providedHash, user.PassHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if the folders exists
	folders, err := c.ci.Read(c.ctx, user)
	if err != nil {
		c.logger.Println(err)
		http.Error(w, "No folders found", http.StatusNotFound)
		return
	}

	for _, v := range folders {
		err = c.ci.Delete(c.ctx, user, v.Name)
		if err != nil {
			c.logger.Println(err)
			http.Error(w, "Error processing Delete", http.StatusInternalServerError)
		}
	}

	c.logger.Printf("Deleted folders for %s\n", user)
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusNoContent)
}

func checkHash(a, b [sha256.Size]byte) bool {
	for i := 0; i < sha256.Size; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
