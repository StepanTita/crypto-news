package ctx

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"common/data/store"
)

const (
	ctxLog          = "ctxLog"
	ctxDataProvider = "ctxDataProvider"
	ctxOAuth        = "ctxOAuth"
)

func CtxLog(log *logrus.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxLog, log)
	}
}

func Log(r *http.Request) *logrus.Entry {
	return r.Context().Value(ctxLog).(*logrus.Entry)
}

func CtxDataProvider(provider store.DataProvider) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxDataProvider, provider)
	}
}

func DataProvider(r *http.Request) store.DataProvider {
	return r.Context().Value(ctxDataProvider).(store.DataProvider)
}

func CtxOAuth(oauth *oauth2.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxOAuth, oauth)
	}
}

func OAuth(r *http.Request) *oauth2.Config {
	return r.Context().Value(ctxOAuth).(*oauth2.Config)
}
