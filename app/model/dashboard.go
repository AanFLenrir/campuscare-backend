package model

type DashboardSummary struct {
	Role        string      `json:"role"`
	Data        interface{} `json:"data"`
	GeneratedAt string      `json:"generated_at"`
}

type AdminDashboardSummary struct {
	TotalUsers           int64            `json:"total_users"`
	TotalClinics         int64            `json:"total_clinics"`
	TotalPets            int64            `json:"total_pets"`
	TotalAppointments    int64            `json:"total_appointments"`
	TotalRecords         int64            `json:"total_medical_records"`
	TotalVaccinations    int64            `json:"total_vaccinations"`
	AppointmentsByStatus map[string]int64 `json:"appointments_by_status"`
}

type PetOwnerDashboardSummary struct {
	TotalMyPets          int64 `json:"total_my_pets"`
	TotalAppointments    int64 `json:"total_appointments"`
	UpcomingAppointments int64 `json:"upcoming_appointments"`
	TotalMedicalRecords  int64 `json:"total_medical_records"`
	TotalVaccinations    int64 `json:"total_vaccinations"`
}

type VeterinarianDashboardSummary struct {
	TotalAssignedAppointments int64 `json:"total_assigned_appointments"`
	TodayAppointments         int64 `json:"today_appointments"`
	CompletedAppointments     int64 `json:"completed_appointments"`
	TotalRecordsCreated       int64 `json:"total_records_created"`
}
