package domains

import (
	"context"
	"go.mercari.io/datastore"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
	"time"
)

type User struct {
	ID                int64  `datastore:"-" boom:"id"`
	Name              string `validate:"req"`
	NewPassword       string `datastore:"-"`
	EncryptedPassword string `validate:"req"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func (user *User) Load(ctx context.Context, ps []datastore.Property) error {
	return datastore.LoadStruct(ctx, user, ps)
}

func (user *User) Save(ctx context.Context) ([]datastore.Property, error) {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	err := user.EncryptIfNeeded()
	if err != nil {
		return nil, err
	}

	err = Validate(user)
	if err != nil {
		return nil, err
	}

	return datastore.SaveStruct(ctx, user)
}

func (user *User) EncryptIfNeeded() error {
	if user.NewPassword == "" {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), 0)
	if err != nil {
		return err
	}
	user.EncryptedPassword = string(hash)
	user.NewPassword = ""
	return nil
}

func (user *User) CheckPassword(input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(input))
	if xerrors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
