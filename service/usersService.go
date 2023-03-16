package service

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/innocentkithinji/xmtest/repository"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type (
	userService struct {
		userRepo repository.UsersRepository
	}

	UserSvc interface {
		Register(user *entity.User) (*entity.LoginCredentials, error)
		Login(user *entity.User) (*entity.LoginCredentials, error)
		Retrieve(uid string) (*entity.User, error)
		Update(uid string, user *entity.User) (*entity.User, error)
		generateAccessCreds(user *entity.User) (*entity.LoginCredentials, error)
	}
)

func (u userService) generateAccessCreds(user *entity.User) (*entity.LoginCredentials, error) {
	log.Printf("User data %+v", user)
	claims := &entity.JWTClaims{
		Email: user.Email,
		UID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	}

	log.Printf("Claims %+v", user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("Token %+v", token)
	signingSecret := viper.Get("SIGNING_SECRET").(string)
	signedToken, err := token.SignedString([]byte(signingSecret))
	log.Printf("Signed %+v", signedToken)
	if err != nil {
		return nil, err
	}

	return &entity.LoginCredentials{Token: signedToken}, nil
}

func (u userService) Register(user *entity.User) (*entity.LoginCredentials, error) {
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Unable to hash password: %s", err)
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.ID = uuid.New().String()
	newUser, err := u.userRepo.Create(user)
	if err != nil {
		log.Printf("Error Creeating new user")
		return nil, err
	}

	return u.generateAccessCreds(newUser)

}

func (u userService) Login(user *entity.User) (*entity.LoginCredentials, error) {
	if user.Email == "" {
		return nil, errors.New("email not provided")
	}
	filter := map[string]interface{}{
		"email": user.Email,
	}
	savedUser, err := u.userRepo.Filter(filter)
	if err != nil {
		log.Println("User with given email was not found")
		return nil, errors.New("Invalid Login Credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("Invalid Login Credentials")
	}
	return u.generateAccessCreds(savedUser)
}

func (u userService) Retrieve(uid string) (*entity.User, error) {
	return u.userRepo.Get(uid)
}

func (u userService) Update(uid string, update *entity.User) (*entity.User, error) {
	user, err := u.userRepo.Get(uid)
	if err != nil {
		log.Printf("Could not find user with ID: %s", uid)
		return nil, err
	}
	if update.Password != "" {
		pwd := update.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Unable to hash password: %s", err)
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if update.Email != "" {
		user.Email = update.Email
	}

	return u.userRepo.Update(user)
}

func NewUsersService(userRepo repository.UsersRepository) UserSvc {
	return userService{userRepo: userRepo}
}
