package todoerrs

import (
	"fmt"

	"github.com/preslavmihaylov/todocheck/checker"
	"github.com/preslavmihaylov/todocheck/fetcher"
	"github.com/preslavmihaylov/todocheck/matchers"
	"github.com/preslavmihaylov/todocheck/traverser/comments"
)

// TodoErrCallback is a function which acts on an encountered todo error
type TodoErrCallback func(todoerr error) error

// NewTraverser for todo errors
func NewTraverser(f *fetcher.Fetcher, ignoredPaths []string, callback TodoErrCallback) *Traverser {
	chk := checker.New(f)
	return &Traverser{
		comments.NewTraverser(ignoredPaths, commentsCallback(chk, callback)),
	}
}

// Traverser for todo errors
type Traverser struct {
	commentsTraverser *comments.Traverser
}

func commentsCallback(chk *checker.Checker, todoErrCallback TodoErrCallback) comments.CommentCallback {
	return func(comment, filepath string, lines []string, linecnt int) error {
		todoErr, err := chk.Check(matchers.ForFile(filepath), comment, filepath, lines, linecnt)
		if err != nil {
			return fmt.Errorf("couldn't check todo line: %w", err)
		} else if todoErr != nil {
			todoErrCallback(todoErr)
		}

		return nil
	}
}

// TraversePath for todo errors. Callback is invoked on encountered error
func (t *Traverser) TraversePath(path string) error {
	return t.commentsTraverser.TraversePath(path)
}