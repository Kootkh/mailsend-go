package engineer

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	engineer "github.com/muquit/mailsend-go/code_examples/api/handler"
	"github.com/samverrall/gopherjobs/api/handler"
	"github.com/samverrall/gopherjobs/api/middleware"
)

func Routes(h *hanldler.Handler) {
	engineerGroup := h.App.Group("/engineer")
	engineerGroup.Use(middleware.RequireAuthentication(h.AccountRepo, h.SB))
	engineerGroup.Get("/dashboard", getDashboard(h))
	engineerGroup.Get("/manage-profile/details", getManageProfile(h))
	engineerGroup.Get("/manage-profile/experience", getManageExperience(h))
	engineerGroup.Get("/add-experience-form", getAddexperienceForm(h))
	engineerGroup.Get("/profile/experience/:experienceID/form", getEditExperienceForm(h))

	apiEngineerGroup := h.App.Group("/api/engineers")
	apiEngineerGroup.Post("/profile", postCreateProfile(h))
	apiEngineerGroup.Post("/profile/experience", postAddProfileExperience(h))
	apiEngineerGroup.Get("/profile/experience/:experienceID", getExperienceItem(h))
	apiEngineerGroup.Patch("/profile/experience/:experienceID", patchEditExperience(h))
}

func getExperienceItem(h *handler.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.Context()

		logger := h.Logger(c)

		// Debug
		// logger.DebugContext(ctx, "get experience item")

		experienceID := c.Params("experienceID")

		if experienceID == "" {
			return h.RenderComponent(c, "edit_experience", errors.New("Missing experience ID"), fiber.Map{
				"Experience": engineer.Experience{},
			})
		}

		exp, err := h.EngineerRepo.GetExperienceByID(ctx, experienceID)

		if err != nil {
			logger.ErrorContext(ctx, "failed to get experience by id", "error", err)

			return h.RenderComponent(c, "edit_experience", err, fiber.Map{
				"Experience": engineer.Experience{},
			})
		}

		return h.RenderComponent(c, "experience_item", nil, fiber.Map{
			"Experience": exp,
		})
	}
}

func patchEditExperience(h *handler.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		experienceID := c.Params("experienceID")

		if experienceID == "" {
			return h.RenderComponent(c, "edit_experience", errors.New("Missing experience ID"), fiber.Map{
				"Experience": engineer.Experience{},
			})
		}

	}
}
