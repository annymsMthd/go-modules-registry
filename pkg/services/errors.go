package services

import (
	"fmt"
)

type ErrModuleDoesntExist struct {
	module string
}

func NewErrModuleDoesntExist(module string) *ErrModuleDoesntExist {
	return &ErrModuleDoesntExist{module}
}

func (e *ErrModuleDoesntExist) Error() string {
	return fmt.Sprintf("module %s does not exit", e.module)
}
