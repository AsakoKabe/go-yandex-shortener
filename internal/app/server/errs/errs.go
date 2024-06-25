package errs

import "fmt"

var ErrConflictOriginalURL = fmt.Errorf("original Url Already Exist")
var ErrCreateDBPoll = fmt.Errorf("error creating db pool")
var ErrCreateServices = fmt.Errorf("error creating db services")
var ErrRegisterEndpoints = fmt.Errorf("error regestration http endpoints")
