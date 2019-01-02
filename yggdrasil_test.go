package l

import (
	"fmt"
	"testing"

	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
)

type testRegister struct {
	name   string
	logger Logger
	err    error
}

func TestRegister(test *testing.T) {
	scenarios := []testRegister{
		testRegister{
			name:   "Register the Logger reference",
			logger: NewZapLoggerDefault(),
		},
		testRegister{
			name:   "Register a nil Logger reference",
			logger: nil,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				roots := yggdrasil.NewRoots()
				err := Register(&roots, scenario.logger)
				assert.Equal(t, scenario.err, err)

				tree := roots.NewTreeDefault()
				logger, err := tree.Reference(loggerPath)

				assert.Nil(t, err, "tree reference error")
				assert.Exactly(t, scenario.logger, logger, "logger reference")
			},
		)
	}
}

type testReference struct {
	name       string
	references map[yggdrasil.Path]yggdrasil.Reference
	tree       yggdrasil.Tree
	err        error
}

func (scenario *testReference) setup(t *testing.T) {
	roots := yggdrasil.NewRoots()
	for path, reference := range scenario.references {
		err := roots.Register(path, reference)
		assert.Nil(t, err, "testReferenceSetup error")
	}
	scenario.tree = roots.NewTreeDefault()
}

func TestReference(test *testing.T) {
	scenarios := []testReference{
		testReference{
			name: "Access the Logger Reference",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				loggerPath: yggdrasil.NewReference(NewZapLoggerDefault()),
			},
		},
		testReference{
			name: "Access a nil Logger Reference",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				loggerPath: yggdrasil.NewReference(nil),
			},
		},
		testReference{
			name: "When Logger was not register returns path not found",
			err:  yggdrasil.ErrPathNotFound,
		},
		testReference{
			name: "When a invalid Logger was register returns invalid reference error",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				loggerPath: yggdrasil.NewReference(new(struct{})),
			},
			err: ErrInvalidReference,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				_, err := Reference(scenario.tree)
				assert.Equal(t, scenario.err, err, "reference error")
				if scenario.err != nil {
					assert.PanicsWithValue(t, scenario.err,
						func() {
							_ = MustReference(scenario.tree)
						},
					)
				} else {
					assert.NotPanics(t,
						func() {
							_ = MustReference(scenario.tree)
						},
					)
				}
			},
		)
	}
}
