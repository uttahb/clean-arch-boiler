package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserInvalidCredentials = errors.New("user invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")
