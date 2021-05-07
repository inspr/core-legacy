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
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/controller/client"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/utils"
)

// Client defines a Redis client, which has the interface methods
type Client struct {
	rdb           *redis.ClusterClient
	refreshURL    string
	refreshKey    string
	insprdAddress string
}

func (c *Client) initAdminUser() error {
	adminUser := User{
		UID:         "admin",
		Permissions: map[string][]string{"": {auth.CreateToken}},
		Password:    os.Getenv("ADMIN_PASSWORD"),
	}
	payload, _ := c.encrypt(adminUser)
	token, err := c.requestNewToken(context.Background(), *payload)
	if err != nil {
		return err
	}
	os.Setenv("ADMIN_TOKEN", token)
	return set(context.Background(), c.rdb, adminUser)
}

// NewRedisClient creates and returns a new Redis client
func NewRedisClient() *Client {
	c := &Client{
		rdb: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{fmt.Sprintf("%s:%s", getEnv("REDIS_HOST"), getEnv("REDIS_PORT"))},
			Password: getEnv("REDIS_PASSWORD"),
		}),
		refreshURL:    fmt.Sprintf("%s/refreshtoken", getEnv("REFRESH_URL")),
		refreshKey:    getEnv("REFRESH_KEY"),
		insprdAddress: getEnv("INSPR_CLUSTER_ADDR"),
	}
	err := c.initAdminUser()
	if err != nil {
		fmt.Println("ERROR CREATING REDIS-CLIENT", err.Error())
		log.Println("ERROR CREATING REDIS-CLIENT", err.Error())
		panic(err)
	}
	return c
}

// CreateUser inserts a new user into Redis
func (c *Client) CreateUser(ctx context.Context, uid, pwd string, newUser User) error {
	err := hasPermission(ctx, c.rdb, uid, pwd)

	if err != nil {
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}
	if err := set(ctx, c.rdb, newUser); err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	return nil
}

// DeleteUser deletes an user from Redis, if it exists
func (c *Client) DeleteUser(ctx context.Context, uid, pwd, usrToBeDeleted string) error {
	err := hasPermission(ctx, c.rdb, uid, pwd)

	if err != nil {
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}
	if err = delete(ctx, c.rdb, usrToBeDeleted); err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	return nil
}

// UpdatePassword changes an users password, if that user exists
func (c *Client) UpdatePassword(ctx context.Context, uid, pwd, usrToBeUpdated, newPwd string) error {
	err := hasPermission(ctx, c.rdb, uid, pwd)

	if err != nil {
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}
	user, err := get(ctx, c.rdb, usrToBeUpdated)
	if err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	user.Password = newPwd

	if err := set(ctx, c.rdb, *user); err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	return nil
}

// Login receives an user and a password, and checks if they exist and match.
// If so, it sends a request to Insprd so it can generate a new token for the
// given user, and returns the toker if it's creation was successful
func (c *Client) Login(ctx context.Context, uid, pwd string) (string, error) {
	user, err := get(ctx, c.rdb, uid)
	if err != nil {
		return "", ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	if pwd != user.Password {
		return "", ierrors.NewError().Unauthorized().
			Message("user and password don't match").Build()
	}

	payload, err := c.encrypt(*user)
	if err != nil {
		return "", ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	token, err := c.requestNewToken(ctx, *payload)
	if err != nil {
		return "", ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	return token, nil
}

// RefreshToken receives a refreshToken and checks if it's valid.
// If so, it returns a payload containing the updated user info
// (user which is associated with the given refreshToken)
func (c *Client) RefreshToken(ctx context.Context, refreshToken []byte) (*auth.Payload, error) {
	oldUser, err := c.decrypt(refreshToken)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	newUser, err := get(ctx, c.rdb, oldUser.UID)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	updatedPayload, err := c.encrypt(*newUser)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	return updatedPayload, nil
}

func (c *Client) encrypt(user User) (*auth.Payload, error) {
	stringToEncrypt := fmt.Sprintf("%s:%s", user.UID, user.Password)

	//Since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(c.refreshKey)
	if err != nil {
		return nil, err
	}
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	payload := auth.Payload{
		UID:         user.UID,
		Permissions: user.Permissions,
		Refresh:     ciphertext,
		RefreshURL:  c.refreshURL,
	}

	return &payload, nil
}

func (c *Client) decrypt(encryptedString []byte) (*User, error) {
	key, err := hex.DecodeString(c.refreshKey)
	if err != nil {
		return nil, err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(encryptedString) < nonceSize {
		return nil, fmt.Errorf("invalid refresh token")
	}

	//Extract the nonce from the encrypted data
	nonce, ciphertext := encryptedString[:nonceSize], encryptedString[nonceSize:]

	//Decrypt the data
	bytetext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	plaintext := string(bytetext)
	usrData := strings.Split(plaintext, ":")
	if len(usrData) != 2 {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return &User{UID: usrData[0], Password: usrData[1]}, nil
}

type authorizer struct{}

func (authorizer) GetToken() ([]byte, error) {
	return []byte("Bearer " + os.Getenv("ADMIN_TOKEN")), nil
}
func (authorizer) SetToken(token []byte) error {
	os.Setenv("ADMIN_TOKEN", string(token)[len("Bearer "):])
	return nil
}

func (c *Client) requestNewToken(ctx context.Context, payload auth.Payload) (string, error) {
	config := client.ControllerConfig{
		Auth: authorizer{},
		URL:  c.insprdAddress,
	}

	ncc := client.NewControllerClient(config)

	token, err := ncc.Authorization().GenerateToken(ctx, payload)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Auxiliar methods

func set(ctx context.Context, rdb *redis.ClusterClient, data User) error {
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

func get(ctx context.Context, rdb *redis.ClusterClient, key string) (*User, error) {
	var parsedValue User
	value, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("key '%v' does not exist", key)
	} else if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(value), &parsedValue)

	return &parsedValue, nil
}

func delete(ctx context.Context, rdb *redis.ClusterClient, key string) error {
	numDeleted, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	} else if numDeleted == 0 {
		return fmt.Errorf("no items were deleted for key %v", key)
	}
	return nil
}

func hasPermission(ctx context.Context, rdb *redis.ClusterClient, uid, pwd string) error {
	requestor, err := get(ctx, rdb, uid)
	if err != nil {
		return err
	}
	if requestor.Password != pwd {
		return fmt.Errorf("invalid password for user %v", uid)
	}

	if rootPerm, ok := requestor.Permissions[""]; ok {
		if utils.Includes(rootPerm, string(auth.CreateToken)) {
			return nil
		}
	}
	return fmt.Errorf("user %v doesn't have admin permission", uid)
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found")
}
