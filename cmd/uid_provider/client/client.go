package client

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/go-redis/redis/v8"
	"gitlab.inspr.dev/inspr/core/pkg/api/auth"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// Client defines a Redis client, which has the interface methods
type Client struct {
	rdb *redis.Client
}

// NewRedisClient creates and returns a new Redis client
func NewRedisClient() *Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	return &Client{
		rdb: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       0, // use default DB
		}),
	}
}

// CreateUser inserts a new user into Redis
func (c *Client) CreateUser(ctx context.Context, uid string, newUser User) error {
	if isAdmin(ctx, c.rdb, uid) {
		if err := set(ctx, c.rdb, newUser); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("current user doesn't have permission to create new users")
}

// DeleteUser deletes an user from Redis, if it exists
func (c *Client) DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error {
	if isAdmin(ctx, c.rdb, uid) {
		if err := delete(ctx, c.rdb, usrToBeDeleted); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("current user doesn't have permission to create new users")
}

// UpdatePassword changes an users password, if that user exists
func (c *Client) UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error {
	if isAdmin(ctx, c.rdb, uid) {

		user, err := get(ctx, c.rdb, usrToBeUpdated)
		if err != nil {
			return err
		}

		user.Password = newPwd

		if err := set(ctx, c.rdb, user); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("current user doesn't have permission to create new users")
}

// Login receives an user and a password, and checks if they exist and match.
// If so, it sends a request to Insprd so it can generate a new token for the
// given user, and returns the toker if it's creation was successful
func (c *Client) Login(ctx context.Context, uid, pwd string) (string, error) {
	user, err := get(ctx, c.rdb, uid)
	if err != nil {
		return "", err
	}
	if pwd != string(user.Password) {
		return "", fmt.Errorf("user and password don't match")
	}

	payload, err := encrypt(user)
	if err != nil {
		return "", err
	}

	token, err := requestNewToken(ctx, payload)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RefreshToken receives a refreshToken and checks if it's valid.
// If so, it returns a payload containing the updated user info
// (user which is associated with the given refreshToken)
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (auth.Payload, error) {
	oldUser, err := decrypt(refreshToken)
	if err != nil {
		return auth.Payload{}, err
	}

	newUser, err := get(ctx, c.rdb, oldUser.UID)
	if err != nil {
		return auth.Payload{}, err
	}

	updatedPayload, err := encrypt(newUser)
	if err != nil {
		return auth.Payload{}, err
	}

	return updatedPayload, nil
}

// Auxiliar methods

func set(ctx context.Context, rdb *redis.Client, data User) error {
	strData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, data.UID, strData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func get(ctx context.Context, rdb *redis.Client, key string) (User, error) {
	var parsedValue User
	value, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return User{}, fmt.Errorf("key '%v' does not exist", key)
	} else if err != nil {
		return User{}, err
	}

	json.Unmarshal([]byte(value), &parsedValue)

	return parsedValue, nil
}

func delete(ctx context.Context, rdb *redis.Client, key string) error {

	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func isAdmin(ctx context.Context, rdb *redis.Client, uid string) bool {
	requestor, err := get(ctx, rdb, uid)
	if err != nil {
		return false
	}
	if requestor.Role != 1 {
		return false
	}
	return true
}

func encrypt(user User) (auth.Payload, error) {
	keyString := "somehow get it from the cluster"
	stringToEncrypt := fmt.Sprintf("%s:%s", user.UID, user.Password)

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return auth.Payload{}, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return auth.Payload{}, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return auth.Payload{}, err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	payload := auth.Payload{
		UID:     user.UID,
		Role:    user.Role,
		Scope:   user.Scope,
		Refresh: string(ciphertext),
	}

	return payload, nil
}

func decrypt(encryptedString string) (User, error) {
	usr := User{}
	keyString := os.Getenv("REFRESH_KEY")
	if keyString == "" {
		return User{}, fmt.Errorf("decryption key is empty")
	}

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return User{}, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return User{}, err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return User{}, err
	}

	if err := json.Unmarshal(plaintext, &usr); err != nil {
		return User{}, err
	}

	return usr, nil
}

func requestNewToken(ctx context.Context, payload auth.Payload) (string, error) {
	url := os.Getenv("INSPR_CLUSTER_ADDR")
	if url == "" {
		panic("[ENV VAR] INSPR_CLUSTER_ADDR not found")
	}

	rc := request.NewClient().
		BaseURL(url).
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Build()

	ncc := client.NewControllerClient(rc)

	token, err := ncc.Authorization().GenerateToken(ctx, payload)
	if err != nil {
		return "", err
	}

	return token, nil
}
