package domain

import (
	"database/sql"

	"github.com/go-hexagonal-arch-auth/errs"
	"github.com/go-hexagonal-arch-auth/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT username, password, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ?
                GROUP BY a.customer_id`

	err := d.client.Get(&login, sqlVerify, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while decrypting password: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected decrypt error")
		}
	}

	return &login, nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate the refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}

	// store it in the store
	sqlInsert := "insert into refresh_token_store (refresh_token) values (?)"
	_, err := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
