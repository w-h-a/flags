package github

import (
	"context"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

type githubTokenKey struct{}

func WithGithubToken(token string) file.Option {
	return func(o *file.Options) {
		o.Context = context.WithValue(o.Context, githubTokenKey{}, token)
	}
}

func GithubToken(ctx context.Context) (string, bool) {
	t, ok := ctx.Value(githubTokenKey{}).(string)
	return t, ok
}
