package engineer

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/samverrall/gopherjobs/api/apierror"
	"github.com/samverrall/gopherjobs/api/handler"
	"github.com/samverrall/gopherjobs/api/templateutil"
	"github.com/samverrall/gopherjobs/internal/app"
	"github.com/samverrall/gopherjobs/internal/app/engineer"
)

func postAddProfileExperience(h *handler.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		user := h.User(c)

		if user == nil {
			return apierror.ErrorResponse(c, apierror.ErrUnauthorized)
		}

		logger := h.Logger(c)
		passport := h.Passport(c)

		var req engineer.AddExperienceInput

		if err := c.BodyParser(&req); err != nil {
			logger.ErrorContext(ctx, "failed to parse body", "error", err)
			return templateutil.Render(h.View, c, "engineer/com_experience_item", fiber.Map{
				"Error": "invalid input provided",
			}, "layouts/none")
		}

		profile, err := h.EngineerRepo.GetProfileForUser(ctx, user.UUID)

		if err != nil {

			logger.ErrorContext(ctx, "failed to get profile for user", "error", err)

			return h.RenderComponent(c, "experience_item", err, fiber.Map{
				"Experience": "Something went wrong",
			})
		}

		exp, err := h.EngineerAPI.AddExperience(ctx, passport.Engineer, profile.UUID, user.UUID, req)

		if err != nil {
			logger.ErrorContext(ctx, "failed to add experience", "error", err)

			if errors.Is(err, app.ErrInvalidInput) {
				return templateutil.Render(h.View, c, "engineer/com_experience_item", fiber.Map{
					"Error": err.Error(),
				}, "layouts/none")
			}

			return templateutil.Render(h.View, c, "engineer/com_experience_item", fiber.Map{
				"Error": "Something went wrong",
			}, "layouts/none")
		}

		return h.RenderComponent(c, "experience_item", nil, fiber.Map{
			"Experience": exp,
		})
	}
}
