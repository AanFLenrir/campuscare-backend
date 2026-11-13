package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
	"github.com/AanFLenrir/campuscare-backend/app/service"
	"github.com/AanFLenrir/campuscare-backend/config"
	"github.com/AanFLenrir/campuscare-backend/helper"
	"github.com/AanFLenrir/campuscare-backend/middleware"
)

// Setup registers all application routes.
func Setup(app *fiber.App, cfg *config.Config, db *pgxpool.Pool) {
	// ── Repositories ──────────────────────────────────────────────────────────
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	complaintRepo := repository.NewComplaintRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	historyRepo := repository.NewComplaintHistoryRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, roleRepo, cfg)
	facilitySvc := service.NewFacilityService(facilityRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	complaintSvc := service.NewComplaintService(complaintRepo, assignmentRepo, historyRepo, userRepo, roleRepo)
	assignmentSvc := service.NewAssignmentService(assignmentRepo)
	userSvc := service.NewUserService(userRepo, roleRepo)
	reportSvc := service.NewReportService(complaintRepo, assignmentRepo, historyRepo)

	// ── Public Routes ──────────────────────────────────────────────────────────
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth
	api.Post("/auth/register", func(c *fiber.Ctx) error {
		var req struct {
			Name          string `json:"name"`
			Email         string `json:"email"`
			Password      string `json:"password"`
			Phone         string `json:"phone"`
			StudentNumber string `json:"student_number"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		result, err := authSvc.Register(context.Background(), service.RegisterInput{
			Name: req.Name, Email: req.Email, Password: req.Password,
			Phone: req.Phone, StudentNumber: req.StudentNumber,
		})
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, result)
	})

	api.Post("/auth/login", func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		result, err := authSvc.Login(context.Background(), req.Email, req.Password)
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, err.Error())
		}
		return helper.Success(c, result)
	})

	// ── Protected Routes ───────────────────────────────────────────────────────
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Profile
	protected.Get("/auth/me", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		user, err := authSvc.GetProfile(context.Background(), uid)
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "user not found")
		}
		return helper.Success(c, user)
	})

	// ── Admin: User Management ──────────────────────────────────────────────────
	adminUsers := protected.Group("/users", middleware.Authorize("ADMIN"))

	adminUsers.Get("", func(c *fiber.Ctx) error {
		users, err := userSvc.ListUsers(context.Background())
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch users")
		}
		return helper.Success(c, users)
	})

	adminUsers.Get("/:id", func(c *fiber.Ctx) error {
		user, err := userSvc.GetUser(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "user not found")
		}
		return helper.Success(c, user)
	})

	adminUsers.Post("", func(c *fiber.Ctx) error {
		var req struct {
			Name          string `json:"name"`
			Email         string `json:"email"`
			Password      string `json:"password"`
			Phone         string `json:"phone"`
			StudentNumber string `json:"student_number"`
			RoleName      string `json:"role"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		user, err := userSvc.CreateUser(context.Background(), service.CreateUserInput{
			Name: req.Name, Email: req.Email, Password: req.Password,
			Phone: req.Phone, StudentNumber: req.StudentNumber, RoleName: req.RoleName,
		})
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, user)
	})

	adminUsers.Put("/:id", func(c *fiber.Ctx) error {
		var req service.UpdateUserInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		user, err := userSvc.UpdateUser(context.Background(), c.Params("id"), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, user)
	})

	adminUsers.Delete("/:id", func(c *fiber.Ctx) error {
		if err := userSvc.DeleteUser(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to delete user")
		}
		return helper.Success(c, fiber.Map{"message": "user deactivated"})
	})

	// ── Facilities ─────────────────────────────────────────────────────────────
	protected.Get("/facilities", func(c *fiber.Ctx) error {
		list, err := facilitySvc.GetAll(context.Background())
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch facilities")
		}
		return helper.Success(c, list)
	})
	protected.Get("/facilities/:id", func(c *fiber.Ctx) error {
		f, err := facilitySvc.GetByID(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "facility not found")
		}
		return helper.Success(c, f)
	})

	adminFacilities := protected.Group("/facilities", middleware.Authorize("ADMIN"))
	adminFacilities.Post("", func(c *fiber.Ctx) error {
		var f model.Facility
		if err := c.BodyParser(&f); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		if err := facilitySvc.Create(context.Background(), &f); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, f)
	})
	adminFacilities.Put("/:id", func(c *fiber.Ctx) error {
		var f model.Facility
		if err := c.BodyParser(&f); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		f.ID = c.Params("id")
		if err := facilitySvc.Update(context.Background(), &f); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, f)
	})
	adminFacilities.Delete("/:id", func(c *fiber.Ctx) error {
		if err := facilitySvc.Delete(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to delete facility")
		}
		return helper.Success(c, fiber.Map{"message": "deleted"})
	})

	// ── Categories ─────────────────────────────────────────────────────────────
	protected.Get("/categories", func(c *fiber.Ctx) error {
		list, err := categorySvc.GetAll(context.Background())
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch categories")
		}
		return helper.Success(c, list)
	})
	protected.Get("/categories/:id", func(c *fiber.Ctx) error {
		cat, err := categorySvc.GetByID(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "category not found")
		}
		return helper.Success(c, cat)
	})

	adminCategories := protected.Group("/categories", middleware.Authorize("ADMIN"))
	adminCategories.Post("", func(c *fiber.Ctx) error {
		var cat model.Category
		if err := c.BodyParser(&cat); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		if err := categorySvc.Create(context.Background(), &cat); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, cat)
	})
	adminCategories.Put("/:id", func(c *fiber.Ctx) error {
		var cat model.Category
		if err := c.BodyParser(&cat); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		cat.ID = c.Params("id")
		if err := categorySvc.Update(context.Background(), &cat); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, cat)
	})
	adminCategories.Delete("/:id", func(c *fiber.Ctx) error {
		if err := categorySvc.Delete(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to delete category")
		}
		return helper.Success(c, fiber.Map{"message": "deleted"})
	})

	// ── Complaints: MAHASISWA ──────────────────────────────────────────────────
	student := protected.Group("/my-complaints", middleware.Authorize("MAHASISWA"))

	student.Post("", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req service.CreateComplaintInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		comp, err := complaintSvc.CreateComplaint(context.Background(), uid, req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, comp)
	})

	student.Get("", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		list, err := complaintSvc.ListMyComplaints(context.Background(), uid)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch complaints")
		}
		return helper.Success(c, list)
	})

	student.Get("/:id", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		comp, err := complaintSvc.GetMyComplaint(context.Background(), uid, c.Params("id"))
		if err != nil {
			if err.Error() == "forbidden" {
				return helper.Error(c, fiber.StatusForbidden, "access denied")
			}
			return helper.Error(c, fiber.StatusNotFound, "complaint not found")
		}
		return helper.Success(c, comp)
	})

	// ── Complaints: ADMIN ──────────────────────────────────────────────────────
	adminComplaints := protected.Group("/complaints", middleware.Authorize("ADMIN"))

	adminComplaints.Get("", func(c *fiber.Ctx) error {
		list, err := complaintSvc.ListAllComplaints(context.Background())
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch complaints")
		}
		return helper.Success(c, list)
	})

	adminComplaints.Get("/:id", func(c *fiber.Ctx) error {
		comp, err := complaintSvc.GetComplaint(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "complaint not found")
		}
		return helper.Success(c, comp)
	})

	adminComplaints.Put("/:id/status", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req struct {
			Status string `json:"status"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		comp, err := complaintSvc.UpdateComplaintStatus(context.Background(), uid, c.Params("id"), req.Status)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, comp)
	})

	adminComplaints.Post("/:id/assign", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req struct {
			OfficerID string `json:"officer_id"`
			Notes     string `json:"notes"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		assign, err := complaintSvc.AssignComplaint(context.Background(), uid, c.Params("id"), req.OfficerID, req.Notes)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, assign)
	})

	// ── Assignments: PETUGAS ──────────────────────────────────────────────────
	officer := protected.Group("/assignments", middleware.Authorize("PETUGAS"))

	officer.Get("", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		list, err := assignmentSvc.ListMyAssignments(context.Background(), uid)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch assignments")
		}
		return helper.Success(c, list)
	})

	officer.Put("/:id/status", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req struct {
			Status string `json:"status"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid payload")
		}
		comp, err := complaintSvc.UpdateComplaintStatus(context.Background(), uid, c.Params("id"), req.Status)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, comp)
	})

	// ── Reports: ADMIN ─────────────────────────────────────────────────────────
	reports := protected.Group("/reports", middleware.Authorize("ADMIN"))

	reports.Get("/summary", func(c *fiber.Ctx) error {
		summary, err := reportSvc.GetSummary(context.Background())
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, "failed to generate report")
		}
		return helper.Success(c, summary)
	})
}
