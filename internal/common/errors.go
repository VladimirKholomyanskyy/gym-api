package common

import "errors"

var ErrAccessForbidden = errors.New("access to resource forbidden for user")
var ErrProgramNotFound = errors.New("training program not found")
