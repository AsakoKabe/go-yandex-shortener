package errs

import "fmt"

// ErrConflictOriginalURL Ошибка для обозначения, что этот URL уже сжали
var ErrConflictOriginalURL = fmt.Errorf("original Url Already Exist")

// ErrCreateDBPoll Ошибка при подключения к БД
var ErrCreateDBPoll = fmt.Errorf("error creating db pool")

// ErrCreateServices Ошибка создания сервисов к БД
var ErrCreateServices = fmt.Errorf("error creating db services")

// ErrRegisterEndpoints Ошибка при регистрации endpoints
var ErrRegisterEndpoints = fmt.Errorf("error regestration http endpoints")
