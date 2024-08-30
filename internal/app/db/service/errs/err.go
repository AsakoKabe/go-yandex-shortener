package errs

import "fmt"

// ErrOriginalURLAlreadyExist Ошибка при сжатии URL. Этот URL уже сжат
var ErrOriginalURLAlreadyExist = fmt.Errorf("original URL Already Exist")
