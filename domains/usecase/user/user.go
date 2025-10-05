package user

import (
	"context"
	"database/sql"
	"kc-ewallet/configurations"
	"kc-ewallet/domains/repository"
	"kc-ewallet/domains/repository/postgres"
	"kc-ewallet/internals/errors"
	log_color "kc-ewallet/internals/helpers/color"
	jwtHelper "kc-ewallet/internals/helpers/jwt"
	strhelper "kc-ewallet/internals/helpers/str"
	"kc-ewallet/protocols/http/request"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

type userUsecase struct {
	db         *sql.DB
	repository repository.IRepository
	jwthelpers configurations.IJWTConfiguration
	tracer     trace.Tracer
}

func NewUserUsecase(
	db *sql.DB,
	repository repository.IRepository,
	jwtConfig configurations.IJWTConfiguration,
	trace trace.Tracer,
) *userUsecase {
	return &userUsecase{
		db:         db,
		repository: repository,
		jwthelpers: jwtConfig,
		tracer:     trace,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, request request.RegisterUserRequest) error {
	passwordHash, err := strhelper.Hash(request.Password)
	if err != nil {
		log_color.PrintRedf("CreateUser failed to hash password: %v\n", err)
		return errors.InternalServer.NewWithUserMsg(err, "failed to create user")
	}

	if _, err := u.repository.CreateUser(ctx, postgres.CreateUserParams{
		Username: request.Username,
		Password: passwordHash,
	}); err != nil {
		log_color.PrintRedf("CreateUser failed to create user: %v\n", err)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.BadRequest.NewWithUserMsg(err, "username already exists")
			}
		}
		return errors.InternalServer.NewWithUserMsg(err, "failed to create user")
	}

	return nil
}

func (u *userUsecase) Login(ctx context.Context, request request.LoginRequest) (string, *postgres.User, error) {
	user, err := u.repository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		log_color.PrintRedf("Login failed to get user by username: %v\n", err)
		if err == sql.ErrNoRows {
			return "", nil, errors.NotFound.NewWithUserMsg(err, "user not found")
		}
		return "", nil, errors.InternalServer.NewWithUserMsg(err, "failed to get user by username")
	}

	if !strhelper.CheckHash(user.Password, request.Password) {
		log_color.PrintRedf("Login password mismatch\n")
		return "", nil, errors.Unauthorized.New("invalid credentials")
	}

	accessTokenEXP := time.Now().Add(time.Duration(u.jwthelpers.GetExpireInMinute()) * time.Minute)
	accessTokenClaims := jwtHelper.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "kc-ewallet",
			ExpiresAt: jwt.NewNumericDate(accessTokenEXP),
		},
		UserID: user.ID,
	}

	// refresh token can be implemented later if needed

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(u.jwthelpers.GetSigningKey()))
	if err != nil {
		log_color.PrintRedf("Login failed to create access token: %v\n", err)
		return "", nil, errors.InternalServer.NewWithUserMsg(err, "failed to login")
	}

	return accessToken, &user, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, userID int32) (*postgres.User, error) {
	user, err := u.repository.GetUserByIDLock(ctx, userID)
	if err != nil {
		log_color.PrintRedf("GetUserByID failed to get user by id: %v\n", err)
		if err == sql.ErrNoRows {
			return nil, errors.NotFound.NewWithUserMsg(err, "user not found")
		}
		return nil, errors.InternalServer.NewWithUserMsg(err, "failed to get user by id")
	}

	return &user, nil
}
