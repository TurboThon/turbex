package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(b64password string) (string, error) {
	salt := []byte("salt")
	password, err := base64.StdEncoding.DecodeString(b64password)
	if err != nil {
		return "", err
	}
	userHash := pbkdf2.Key([]byte(password), salt, 600000, 128, sha256.New)
	b64salt := base64.StdEncoding.EncodeToString(salt)
	b64hash := base64.StdEncoding.EncodeToString(userHash)
	return fmt.Sprintf("pbkdf2_sha256:600000:%s:%s", b64salt, b64hash), nil
}

// Spec:
// 128 bits salt -> 16 bytes (ANSSI)
// 600_000 rounds (from OWASP 2023)
// HMAC-SHA-256 (FIPS)
// Using hashcat formet : "pbkdf2_sha256", ":", iterations, ":", base64 salt, ":", base64 digest

// For some reason the pbkdf2_hmac_sha256 rust crate computes pbkdf2_sha256 instead...
//! From pbkdf2_hmac rust crate
//! use hex_literal::hex;
//! use pbkdf2::{pbkdf2_hmac, pbkdf2_hmac_array};
//! use sha2::Sha256;
//!
//! let password = b"password";
//! let salt = b"salt";
//! // number of iterations
//! let n = 600_000;
//! // Expected value of generated key
//! let expected = hex!("669cfe52482116fda1aa2cbe409b2f56c8e45637");

func VerifyPassword(b64password string, dbHash string) bool {
	// Example: pbkdf2_sha256:1000:MTc3MTA0MTQwMjQxNzY=:PYjCU215Mi57AYPKva9j7mvF4Rc5bCnt
	parts := strings.Split(dbHash, ":")
	if len(parts) != 4 {
		log.Println("Error while parsing stored password: failed to find 4 fields")
		return false
	}
	numberOfRounds, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Println("Failed to convert to integer during hash parsing")
		return false
	}
	passwordSpec := struct {
		primitive  string
		workFactor int
		b64salt    string
		b64hash    string
	}{
		primitive:  parts[0],
		workFactor: numberOfRounds,
		b64salt:    parts[2],
		b64hash:    parts[3],
	}
	if passwordSpec.primitive != "pbkdf2_sha256" {
		log.Println("Unsupported hashing primitive")
		return false
	}
	bytesSalt, err := base64.StdEncoding.DecodeString(passwordSpec.b64salt)
	if err != nil {
		log.Println(err)
		return false
	}
	bytesHash, err := base64.StdEncoding.DecodeString(passwordSpec.b64hash)
	if err != nil {
		log.Println(err)
		return false
	}
	bytesPassword, err := base64.StdEncoding.DecodeString(b64password)
	if err != nil {
		log.Println(err)
		return false
	}
	userHash := pbkdf2.Key(bytesPassword, bytesSalt, passwordSpec.workFactor, 128, sha256.New)
	return bytes.Equal(userHash, bytesHash)
}

func isValidLogin(db *mongo.Database, username string, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbLoginRequest models.APILoginRequest
	err := db.Collection(consts.COLLECTION_USER).FindOne(ctx, bson.M{"username": username}).Decode(&dbLoginRequest)
	if err != nil {
		return false
	}
	return VerifyPassword(password, dbLoginRequest.Password)
}

func createSessionCookie(db *mongo.Database, userName string, duration time.Duration) (string, error) {
	randomBytes := make([]byte, 40)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(randomBytes)[:20]

	session := models.Session{
		Id:             primitive.NewObjectID(),
		CookieValue:    token,
		UserName:       userName,
		ExpirationDate: time.Now().Add(duration).UTC().Format(consts.DATE_FORMAT),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection(consts.COLLECTION_SESSIONS).InsertOne(ctx, session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func DoLogin(c *gin.Context, db *mongo.Database, env *structs.Env) {
	var loginData models.APILoginRequest

	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isValid := isValidLogin(db, loginData.UserName, loginData.Password)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong login data"})
		return
	}

	var user models.APIUserDetails
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Collection(consts.COLLECTION_USER).FindOne(ctx, bson.M{"username": loginData.UserName}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	cookie, err := createSessionCookie(db, user.UserName, time.Duration(env.SessionDurationSeconds)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.SetCookie(consts.SESSION_COOKIE_NAME, cookie, env.SessionDurationSeconds, "/", c.GetHeader(consts.HOST_HEADER), true, true)
	c.JSON(http.StatusOK, user)
}
