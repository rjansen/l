package l

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Logger Reference")
	loggerPath          = yggdrasil.NewPath("/l/logger")
)

func Register(roots *yggdrasil.Roots, logger Logger) error {
	return roots.Register(loggerPath, logger)
}

func Reference(tree yggdrasil.Tree) (Logger, error) {
	reference, err := tree.Reference(loggerPath)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, nil
	}
	options, is := reference.(Logger)
	if !is {
		return nil, ErrInvalidReference
	}
	return options, nil
}

func MustReference(tree yggdrasil.Tree) Logger {
	logger, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return logger
}
