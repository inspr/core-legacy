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
	"strings"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/controller/client"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/utils"
)

var logger *zap.Logger

func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "redis-client")))
}

// Client defines a Redis client, which has the interface methods
type Client struct {
	rdb           *redis.ClusterClient
	refreshURL    string
	refreshKey    string
	insprdAddress string
}

func (c *Client) initAdminUser() error {

	password := os.Getenv("")
	if password == "" {
		logger.Error("password is empty")
		return ierrors.NewError().InvalidArgs().Message("invalid password").Build()
	}
	logger.Debug("received password, generating encryption")
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	adminUser := User{
		UID:         "admin",
		Permissions: auth.AdminPermissions,
		Password:    string(hashedPwd),
	}

	logger.Debug("encrypting admin user")
	payload, err := c.encrypt(adminUser)
	if err != nil {
		logger.Error("error encrypting admin user", zap.Error(err))
		return err
	}
	logger.Info("requesting new token from insprd")
	token, err := c.requestNewToken(context.Background(), *payload)
	if err != nil {
		logger.Error("error requesting new token", zap.Error(err), zap.String("insprd-address", c.insprdAddress))
		return err
	}
	os.Setenv("ADMIN_TOKEN", token)
	return set(context.Background(), c.rdb, adminUser)
}

// NewRedisClient creates and returns a new Redis client
func NewRedisClient() *Client {
	password := getEnv("REDIS_PASSWORD")
	refreshURL := getEnv("REFRESH_URL")
	refreshKey := getEnv("REFRESH_KEY")
	insprdAddress := getEnv("INSPR_CLUSTER_ADDRESS")
	redisHost := getEnv("REDIS_HOST")
	redisPort := getEnv("REDIS_PORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
	refreshPath := fmt.Sprintf("%s/refreshtoken", refreshURL)
	logger.Info("initializing redis client", zap.String("redis-address", redisAddress), zap.String("refresh-path", refreshPath), zap.String("insprd-address", insprdAddress))
	c := &Client{
		rdb: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{redisAddress},
			Password: password,
		}),
		refreshURL:    refreshPath,
		refreshKey:    refreshKey,
		insprdAddress: insprdAddress,
	}

	err := c.initAdminUser()
	if err != nil {
		logger.Panic("error initializing admin user", zap.Error(err))
	}
	return c
}

// CreateUser inserts a new user into Redis
func (c *Client) CreateUser(ctx context.Context, uid, pwd string, newUser User) error {
	logger.Debug("checking user permissions", zap.String("user", uid))
	err := hasPermission(ctx, c.rdb, uid, pwd, newUser, true)

	if err != nil {
		logger.Info("unable to aquire permissions", zap.String("user", uid))
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("unable to hash password", zap.Error(err))
		return ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	newUser.Password = string(hashedPwd)

	if err := set(ctx, c.rdb, newUser); err != nil {
		logger.Error("unable to set key on redis", zap.Error(err))
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	return nil
}

// DeleteUser deletes an user from Redis, if it exists
func (c *Client) DeleteUser(ctx context.Context, uid, pwd, usrToBeDeleted string) error {
	user, err := get(ctx, c.rdb, usrToBeDeleted)
	if err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	err = hasPermission(ctx, c.rdb, uid, pwd, *user, false)
	if err != nil {
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}

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
	user, err := get(ctx, c.rdb, usrToBeUpdated)
	if err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	err = hasPermission(ctx, c.rdb, uid, pwd, *user, false)
	if err != nil {
		return ierrors.NewError().Forbidden().Message(err.Error()).Build()
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	user.Password = string(hashedPwd)

	if err := set(ctx, c.rdb, *user); err != nil {
		return ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}
	return nil
}

// Login receives an user and a password, and checks if they exist and match.
// If so, it sends a request to Insprd so it can generate a new token for the
// given user, and returns the toker if it's creation was successful
func (c *Client) Login(ctx context.Context, uid, pwd string) (string, error) {
	logger.Debug("getting user key from redis", zap.String("user", uid))
	user, err := get(ctx, c.rdb, uid)
	if err != nil {
		logger.Error("unable to get key from redis", zap.String("user", uid), zap.Error(err))
		return "", ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	logger.Debug("comparing password")
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)); err != nil {
		logger.Debug("passwords don't match", zap.Error(err))
		return "", ierrors.NewError().Unauthorized().
			Message("user and password don't match").Build()
	}

	logger.Debug("encrypting user")
	payload, err := c.encrypt(*user)
	if err != nil {
		logger.Error("unable to encrypt user", zap.Any("user", user), zap.Error(err))
		return "", ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	logger.Debug("requesting user token")
	token, err := c.requestNewToken(ctx, *payload)
	if err != nil {
		logger.Error("unable to request token from insprd", zap.String("insprd-address", c.insprdAddress), zap.Error(err))
		return "", ierrors.NewError().InternalServer().Message(err.Error()).Build()
	}

	return token, nil
}

// RefreshToken receives a refreshToken and checks if it's valid.
// If so, it returns a payload containing the updated user info
// (user which is associated with the given refreshToken)
func (c *Client) RefreshToken(ctx context.Context, refreshToken []byte) (*auth.Payload, error) {
	logger.Debug("decripting refresh token")
	oldUser, err := c.decrypt(refreshToken)
	if err != nil {
		logger.Error("unable do decript token", zap.Error(err))
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	logger.Debug("retrieving user from redis")
	newUser, err := get(ctx, c.rdb, oldUser.UID)
	if err != nil {
		logger.Error("unable to get key from redis", zap.Any("user", oldUser), zap.Error(err))
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	logger.Debug("encripting new user")
	updatedPayload, err := c.encrypt(*newUser)
	if err != nil {
		logger.Error("unable to encrypt new user", zap.Any("user", newUser), zap.Error(err))
		return nil, ierrors.NewError().BadRequest().Message(err.Error()).Build()
	}

	return updatedPayload, nil
}

func (c *Client) encrypt(user User) (*auth.Payload, error) {
	stringToEncrypt := fmt.Sprintf("%s:%s", user.UID, user.Password)

	//Since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(c.refreshKey)
	if err != nil {
		logger.Error("unable to decode refresh key", zap.Error(err))
		return nil, err
	}
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("unable to initiate cypher from key", zap.Error(err))
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
		logger.Error("error decoding hex stream", zap.Error(err))
		return nil, err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("unable to create aes cipher from key", zap.Binary("key", key), zap.Error(err))
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("invalid GCM block", zap.Error(err))
		return nil, err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(encryptedString) < nonceSize {
		logger.Error("invalid refresh token", zap.String("error", "invalid nonce"))
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
		logger.Error("invalid refresh token", zap.String("error", "not enough fields on refresh token"))
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
		logger.Error("error generating new token", zap.String("insprd-address", c.insprdAddress), zap.Error(err))
		return "", err
	}

	return token, nil
}

// Auxiliar methods

func set(ctx context.Context, rdb *redis.ClusterClient, data User) error {
	strData, err := json.Marshal(data)
	if err != nil {
		logger.Error("error marshalling data", zap.Any("data", data), zap.Error(err))
		return err
	}

	logger.Debug("setting key on redis")
	err = rdb.Set(ctx, data.UID, strData, 0).Err()
	if err != nil {
		logger.Error("error setting key on redis", zap.String("key", data.UID), zap.Error(err))
		return err
	}
	return nil
}

func get(ctx context.Context, rdb *redis.ClusterClient, key string) (*User, error) {
	var parsedValue User
	logger.Debug("retrieving key from redis")
	value, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("key '%v' does not exist", key)
	} else if err != nil {
		logger.Error("error retrieving key from redis", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &parsedValue)
	if err != nil {
		logger.Error("error unmarshalling value from redis", zap.String("value", value), zap.Error(err))
		return nil, err
	}

	return &parsedValue, nil
}

func delete(ctx context.Context, rdb *redis.ClusterClient, key string) error {
	numDeleted, err := rdb.Del(ctx, key).Result()
	if err != nil {
		logger.Error("error deleting redis key", zap.Error(err))
		return err
	} else if numDeleted == 0 {
		logger.Error("key not found")
		return fmt.Errorf("no items were deleted for key %v", key)
	}
	return nil
}

func hasPermission(ctx context.Context, rdb *redis.ClusterClient, uid, pwd string, newUser User, isCreation bool) error {
	requestor, err := get(ctx, rdb, uid)
	if err != nil {
		logger.Error("error getting user", zap.String("user", uid), zap.Error(err))
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(requestor.Password), []byte(pwd)); err != nil {
		logger.Info("invalid password", zap.String("user", uid))
		return fmt.Errorf("invalid password for user %v", uid)
	}

	for newUserPermissionScope, newUserPermissions := range newUser.Permissions {
		isAllowed := false
		for requestorPermissionScope, requestorPermissions := range requestor.Permissions {
			if isPermissionAllowed(newUserPermissionScope, requestorPermissionScope, newUserPermissions, requestorPermissions, isCreation) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			logger.Info("user unauthorized", zap.String("user", uid))
			return ierrors.NewError().Forbidden().Message("not allowed to create/delete/update a user with current permissions").Build()
		}

	}

	return nil
}

func isPermissionAllowed(newUserPermissionScope, requestorPermissionScope string, newUserPermissions, requestorPermissions []string, isCreation bool) bool {
	if !metautils.IsInnerScope(requestorPermissionScope, newUserPermissionScope) {
		return false
	}

	for _, permission := range newUserPermissions {
		if (isCreation && !utils.Includes(requestorPermissions, permission)) || !utils.Includes(requestorPermissions, auth.CreateToken) {
			return false
		}
	}

	return true
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found")
}
