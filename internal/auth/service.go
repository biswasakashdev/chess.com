package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	users "github.com/biswasakashdev/chess.com/internal/models"
	userRepo "github.com/biswasakashdev/chess.com/internal/repository/users"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrUsernameInUse      = errors.New("username already in use")
	ErrInternalError      = errors.New("Somthing went wrong")
)

// AuthService provides authentication functionality
type AuthService struct {
	userRepo       userRepo.UserRepository
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo userRepo.UserRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      []byte(jwtSecret),
		accessTokenTTL: accessTokenTTL,
	}
}

// Register creates a new user with the provided credentials
func (s *AuthService) Register(ctx context.Context, username, password, firstName, lastName string) (*users.User, error) {

	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	// Check if user already exists
	isExist, err := s.userRepo.IsUsernameExits(newCtx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if isExist {
		return nil, ErrUsernameInUse
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user, err := s.userRepo.CreateUser(username, hashedPassword, firstName, lastName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns an access token
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {

	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	// Get the user from the database
	user, err := s.userRepo.FindByUsername(newCtx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Verify the password
	if err := VerifyPassword(user.HashedPassword, password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate an access token
	token, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateAccessToken creates a new JWT access token
func (s *AuthService) generateAccessToken(user *users.User) (string, error) {
	// Set the expiration time
	expirationTime := time.Now().Add(s.accessTokenTTL)

	// Create the JWT claims
	claims := jwt.MapClaims{
		"sub":      user.Id.String(),      // subject (user ID)
		"username": user.Username,         // custom claim
		"exp":      expirationTime.Unix(), // expiration time
		"iat":      time.Now().Unix(),     // issued at time
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
