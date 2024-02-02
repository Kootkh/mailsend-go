package handler

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/samverrall/gopherjobs/api/templateutil"
	"github.com/samverrall/gopherjobs/internal/app/account"
	"github.com/samverrall/gopherjobs/internal/app/engineer"
	"github.com/samverrall/gopherjobs/internal/app/recruiter"
	"github.com/samverrall/gopherjobs/internal/payment"
	"github.com/samverrall/gopherjobs/pkg/config"
)

type Handler struct {
	App             *fiber.App
	View            *templateutil.View
	RecruiterAPI    *recruiter.Service
	EngineerAPI     *engineer.Service
	AccountAPI      *account.Service
	EngineerRepo    engineer.Reader
	RecruiterRepo   recruiter.Reader
	AccountRepo     account.Reader
	NamedEndpoints  templateutil.NamedEndpoints
	SB              *supabase.Client
	Template        *template.Template
	GithubConf      *config.GithubKeys
	PaymentProvider payment.Provider
	Cfg             *config.Config
}

// GithubOAuthConfig - конфигурация OAuth2 для GitHub основанная на учётных данных GitHub получаемых из конфигурации
func (h *Handler) GithubOAuthConfig(ctx context.Context) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.GithubConf.ClientID,
		ClientSecret: h.GithubConf.ClientSecret,
		Scopes:       []string{"user:read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: h.GithubConf.RedirectURL,
	}
}

func (h *Handler) User(c *fiber.Ctx) *account.User {
	u := c.Locals(contextutil.LocalUserKey)

	if user, ok := u.(*account.User); ok {
		return user
	}

	return nil
}

func (h *Handler) Passport(c *fiber.Ctx) *guard.Passport {
	val := c.Locals(contextutil.CtxPassport)

	if val == nil {
		return &guard.Passport{}
	}

	passport, ok := val.(*guard.Passport)

	if !ok {
		panic(fmt.Sprintf("could not assert passport as %T", passport))
	}

	return passport
}

func (h *Handler) Logger(c *fiber.Ctx) *slog.Logger {
	value := c.Locals(contextutil.LoggerCtx{})

	if value == nil {
		return slog.Default()
	}

	logger, ok := value.(*slog.Logger)

	if !ok {
		panic(fmt.Sprintf("could not assert logger as %T", logger))

	}
	return logger
}

func (h *Handler) ExecComponent(c *fiber.Ctx, template string, data fiber.Map) (string, error) {

	logger := h.Logger(c)

	bytesW := new(bytes.Buffer)

	if err := h.Template.ExecuteTemplate(bytesW, template, data); err != nil {
		logger.ErrorContext(c.Context(), "failed to execute template", "error", err)
		return "", apierror.ErrorResponse(c, apierror.ErrorInternal)
	}

	return bytesW.String(), nil
}
