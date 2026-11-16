package route

import (
	"context"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/helper"
	"github.com/AanFLenrir/PetrivaBackend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Setup registers all application routes, repositories, services, and Swagger UI.
// Setup registers all application routes, repositories, services, and Swagger UI.
func Setup(app *fiber.App, cfg *config.Config, db *pgxpool.Pool) {

	// ── Repositories ──────────────────────────────────────────────────────────
	roleRepo := repository.NewRoleRepository(db)
	clinicRepo := repository.NewClinicRepository(db)
	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	apptRepo := repository.NewAppointmentRepository(db)
	recordRepo := repository.NewMedicalRecordRepository(db)
	vacRepo := repository.NewVaccinationRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, roleRepo, cfg)
	userSvc := service.NewUserService(userRepo, roleRepo, clinicRepo)
	clinicSvc := service.NewClinicService(clinicRepo)
	serviceSvc := service.NewServiceService(serviceRepo, clinicRepo)
	petSvc := service.NewPetService(petRepo)
	apptSvc := service.NewAppointmentService(apptRepo, petRepo, userRepo, clinicRepo, serviceRepo)
	recordSvc := service.NewMedicalRecordService(recordRepo, apptRepo, petRepo)
	vacSvc := service.NewVaccinationService(vacRepo, petRepo)
	dashboardSvc := service.NewDashboardService(userRepo, clinicRepo, petRepo, apptRepo, recordRepo, vacRepo)

	// ── Swagger UI Route ──────────────────────────────────────────────────────
	app.Get("/swagger/*", swagger.HandlerDefault)

	// ── API Root & Health ─────────────────────────────────────────────────────
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.Map{"status": "ok", "system": "PETRIVA - Pet Care Management System"})
	})

	// ── Authentication (Public) ───────────────────────────────────────────────
	authGroup := api.Group("/auth")

	authGroup.Post("/register", func(c *fiber.Ctx) error {
		var req model.RegisterInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		res, err := authSvc.Register(context.Background(), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, res)
	})

	authGroup.Post("/login", func(c *fiber.Ctx) error {
		var req model.LoginInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		res, err := authSvc.Login(context.Background(), req)
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, err.Error())
		}
		return helper.Success(c, res)
	})

	// ── Public Catalog: Clinics, Services, Veterinarians ──────────────────────
	api.Get("/clinics", func(c *fiber.Ctx) error {
		clinics, err := clinicSvc.ListClinics(context.Background(), true)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, clinics)
	})

	api.Get("/clinics/:id", func(c *fiber.Ctx) error {
		clinic, err := clinicSvc.GetClinic(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "clinic not found")
		}
		return helper.Success(c, clinic)
	})

	api.Get("/services", func(c *fiber.Ctx) error {
		var clinicID *string
		if q := c.Query("clinic_id"); q != "" {
			clinicID = &q
		}
		services, err := serviceSvc.ListServices(context.Background(), clinicID, true)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, services)
	})

	api.Get("/services/:id", func(c *fiber.Ctx) error {
		svc, err := serviceSvc.GetService(context.Background(), c.Params("id"))
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "service not found")
		}
		return helper.Success(c, svc)
	})

	api.Get("/veterinarians", func(c *fiber.Ctx) error {
		var clinicID *string
		if q := c.Query("clinic_id"); q != "" {
			clinicID = &q
		}
		vets, err := userSvc.ListVeterinarians(context.Background(), clinicID)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, vets)
	})

	api.Get("/veterinarians/:id", func(c *fiber.Ctx) error {
		vet, err := userSvc.GetUser(context.Background(), c.Params("id"))
		if err != nil || vet.Role != model.RoleVeterinarian {
			return helper.Error(c, fiber.StatusNotFound, "veterinarian not found")
		}
		return helper.Success(c, vet)
	})

	// ── Protected Routes (Requires JWT Bearer Token) ──────────────────────────
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Profile
	protected.Get("/auth/me", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		profile, err := authSvc.GetProfile(context.Background(), uid)
		if err != nil {
			return helper.Error(c, fiber.StatusNotFound, "profile not found")
		}
		return helper.Success(c, profile)
	})

	// Dashboard Summary (Role-based)
	protected.Get("/dashboard/summary", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		summary, err := dashboardSvc.GetDashboardSummary(context.Background(), uid, role)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, summary)
	})

	// ── Pets Management (PET_OWNER & ADMIN) ───────────────────────────────────
	protected.Post("/pets", middleware.Authorize(model.RolePetOwner), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req model.CreatePetInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		pet, err := petSvc.CreatePet(context.Background(), uid, req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, pet)
	})

	protected.Get("/pets", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var ownerFilter *string
		if q := c.Query("owner_id"); q != "" && role == model.RoleAdmin {
			ownerFilter = &q
		}
		pets, err := petSvc.ListPets(context.Background(), uid, role, ownerFilter)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, pets)
	})

	protected.Get("/pets/:id", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		pet, err := petSvc.GetPet(context.Background(), uid, role, c.Params("id"))
		if err != nil {
			if err == service.ErrPetForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusNotFound, "pet not found")
		}
		return helper.Success(c, pet)
	})

	protected.Put("/pets/:id", middleware.Authorize(model.RolePetOwner), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req model.UpdatePetInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		pet, err := petSvc.UpdatePet(context.Background(), uid, c.Params("id"), req)
		if err != nil {
			if err == service.ErrPetForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, pet)
	})

	protected.Delete("/pets/:id", middleware.Authorize(model.RolePetOwner), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		if err := petSvc.DeletePet(context.Background(), uid, c.Params("id")); err != nil {
			if err == service.ErrPetForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, fiber.Map{"message": "pet deactivated successfully"})
	})

	// ── Appointments Management ───────────────────────────────────────────────
	protected.Post("/appointments", middleware.Authorize(model.RolePetOwner), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req model.CreateAppointmentInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		appt, err := apptSvc.CreateAppointment(context.Background(), uid, req)
		if err != nil {
			if err == service.ErrPetForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, appt)
	})

	protected.Get("/appointments", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var statusFilter *string
		if q := c.Query("status"); q != "" {
			statusFilter = &q
		}
		appts, err := apptSvc.ListAppointments(context.Background(), uid, role, statusFilter)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, appts)
	})

	protected.Get("/appointments/:id", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		appt, err := apptSvc.GetAppointment(context.Background(), uid, role, c.Params("id"))
		if err != nil {
			if err == service.ErrAppointmentForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusNotFound, "appointment not found")
		}
		return helper.Success(c, appt)
	})

	protected.Put("/appointments/:id/cancel", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var req model.CancelAppointmentInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		appt, err := apptSvc.CancelAppointment(context.Background(), uid, role, c.Params("id"), req.Reason)
		if err != nil {
			if err == service.ErrAppointmentForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, appt)
	})

	protected.Put("/appointments/:id/status", middleware.Authorize(model.RoleVeterinarian, model.RoleAdmin), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var req model.UpdateAppointmentStatusInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		var reason *string
		if req.Reason != "" {
			reason = &req.Reason
		}
		appt, err := apptSvc.UpdateStatus(context.Background(), uid, role, c.Params("id"), req.Status, reason)
		if err != nil {
			if err == service.ErrAppointmentForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, appt)
	})

	// ── Medical Records Management ────────────────────────────────────────────
	protected.Post("/medical-records", middleware.Authorize(model.RoleVeterinarian), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req model.CreateMedicalRecordInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		rec, err := recordSvc.CreateMedicalRecord(context.Background(), uid, req)
		if err != nil {
			if err == service.ErrMedicalRecordForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, rec)
	})

	protected.Get("/medical-records", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var petID *string
		if q := c.Query("pet_id"); q != "" {
			petID = &q
		}
		records, err := recordSvc.ListMedicalRecords(context.Background(), uid, role, petID)
		if err != nil {
			if err == service.ErrMedicalRecordForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, records)
	})

	protected.Get("/medical-records/:id", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		rec, err := recordSvc.GetMedicalRecord(context.Background(), uid, role, c.Params("id"))
		if err != nil {
			if err == service.ErrMedicalRecordForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusNotFound, "medical record not found")
		}
		return helper.Success(c, rec)
	})

	protected.Put("/medical-records/:id", middleware.Authorize(model.RoleVeterinarian, model.RoleAdmin), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var req model.UpdateMedicalRecordInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		rec, err := recordSvc.UpdateMedicalRecord(context.Background(), uid, role, c.Params("id"), req)
		if err != nil {
			if err == service.ErrMedicalRecordForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, rec)
	})

	// ── Vaccinations Management ───────────────────────────────────────────────
	protected.Post("/vaccinations", middleware.Authorize(model.RoleVeterinarian), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		var req model.CreateVaccinationInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		vac, err := vacSvc.CreateVaccination(context.Background(), uid, req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, vac)
	})

	protected.Get("/vaccinations", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var petID *string
		if q := c.Query("pet_id"); q != "" {
			petID = &q
		}
		list, err := vacSvc.ListVaccinations(context.Background(), uid, role, petID)
		if err != nil {
			if err == service.ErrVaccinationForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, list)
	})

	protected.Get("/vaccinations/:id", func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		v, err := vacSvc.GetVaccination(context.Background(), uid, role, c.Params("id"))
		if err != nil {
			if err == service.ErrVaccinationForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusNotFound, "vaccination record not found")
		}
		return helper.Success(c, v)
	})

	protected.Put("/vaccinations/:id", middleware.Authorize(model.RoleVeterinarian, model.RoleAdmin), func(c *fiber.Ctx) error {
		uid, _ := c.Locals("userID").(string)
		role, _ := c.Locals("role").(string)
		var req model.UpdateVaccinationInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		v, err := vacSvc.UpdateVaccination(context.Background(), uid, role, c.Params("id"), req)
		if err != nil {
			if err == service.ErrVaccinationForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, v)
	})

	protected.Delete("/vaccinations/:id", middleware.Authorize(model.RoleVeterinarian, model.RoleAdmin), func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if err := vacSvc.DeleteVaccination(context.Background(), role, c.Params("id")); err != nil {
			if err == service.ErrVaccinationForbidden {
				return helper.Error(c, fiber.StatusForbidden, err.Error())
			}
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, fiber.Map{"message": "vaccination record deleted successfully"})
	})

	// ── Admin: User Management ────────────────────────────────────────────────
	adminUsers := protected.Group("/users", middleware.Authorize(model.RoleAdmin))

	adminUsers.Get("", func(c *fiber.Ctx) error {
		roleFilter := c.Query("role")
		users, err := userSvc.ListUsers(context.Background(), roleFilter)
		if err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
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
		var req model.CreateUserInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		user, err := userSvc.CreateUser(context.Background(), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, user)
	})

	adminUsers.Put("/:id", func(c *fiber.Ctx) error {
		var req model.UpdateUserInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		user, err := userSvc.UpdateUser(context.Background(), c.Params("id"), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, user)
	})

	adminUsers.Delete("/:id", func(c *fiber.Ctx) error {
		if err := userSvc.DeleteUser(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, fiber.Map{"message": "user deactivated successfully"})
	})

	// ── Admin: Clinic Management (POST, PUT, DELETE) ──────────────────────────
	adminClinics := protected.Group("/clinics", middleware.Authorize(model.RoleAdmin))

	adminClinics.Post("", func(c *fiber.Ctx) error {
		var req model.CreateClinicInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		clinic, err := clinicSvc.CreateClinic(context.Background(), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, clinic)
	})

	adminClinics.Put("/:id", func(c *fiber.Ctx) error {
		var req model.UpdateClinicInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		clinic, err := clinicSvc.UpdateClinic(context.Background(), c.Params("id"), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, clinic)
	})

	adminClinics.Delete("/:id", func(c *fiber.Ctx) error {
		if err := clinicSvc.DeleteClinic(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, fiber.Map{"message": "clinic deactivated successfully"})
	})

	// ── Admin: Service Management (POST, PUT, DELETE) ─────────────────────────
	adminServices := protected.Group("/services", middleware.Authorize(model.RoleAdmin))

	adminServices.Post("", func(c *fiber.Ctx) error {
		var req model.CreateServiceInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		svc, err := serviceSvc.CreateService(context.Background(), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Created(c, svc)
	})

	adminServices.Put("/:id", func(c *fiber.Ctx) error {
		var req model.UpdateServiceInput
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
		}
		svc, err := serviceSvc.UpdateService(context.Background(), c.Params("id"), req)
		if err != nil {
			return helper.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return helper.Success(c, svc)
	})

	adminServices.Delete("/:id", func(c *fiber.Ctx) error {
		if err := serviceSvc.DeleteService(context.Background(), c.Params("id")); err != nil {
			return helper.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		return helper.Success(c, fiber.Map{"message": "service deactivated successfully"})
	})
}
