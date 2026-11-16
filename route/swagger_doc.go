package route

import (
	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/helper"
)

var (
	_ = model.RegisterInput{}
	_ = model.LoginInput{}
	_ = model.AuthResponse{}
	_ = model.UserResponse{}
	_ = model.CreateUserInput{}
	_ = model.UpdateUserInput{}
	_ = model.Clinic{}
	_ = model.CreateClinicInput{}
	_ = model.UpdateClinicInput{}
	_ = model.Service{}
	_ = model.CreateServiceInput{}
	_ = model.UpdateServiceInput{}
	_ = model.Pet{}
	_ = model.CreatePetInput{}
	_ = model.UpdatePetInput{}
	_ = model.Appointment{}
	_ = model.CreateAppointmentInput{}
	_ = model.UpdateAppointmentStatusInput{}
	_ = model.CancelAppointmentInput{}
	_ = model.MedicalRecord{}
	_ = model.CreateMedicalRecordInput{}
	_ = model.UpdateMedicalRecordInput{}
	_ = model.Vaccination{}
	_ = model.CreateVaccinationInput{}
	_ = model.UpdateVaccinationInput{}
	_ = model.DashboardSummary{}
	_ = helper.APIResponse{}
)

// ============================================================================
// 1. HEALTH
// ============================================================================

// docHealthCheck godoc
// @Summary Health Check
// @Description Endpoint untuk memeriksa status backend dan ketersediaan layanan
// @Tags Health
// @Produce json
// @Success 200 {object} helper.APIResponse
// @Router /health [get]
func docHealthCheck() {}

// ============================================================================
// 2. AUTHENTICATION
// ============================================================================

// docRegister godoc
// @Summary Register Pet Owner
// @Description Mendaftarkan akun publik baru dengan role PET_OWNER
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterInput true "Payload Registrasi"
// @Success 201 {object} helper.APIResponse{data=model.AuthResponse}
// @Failure 400 {object} helper.APIResponse
// @Router /auth/register [post]
func docRegister() {}

// docLogin godoc
// @Summary User Login
// @Description Otentikasi pengguna menggunakan email dan password untuk mendapatkan JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginInput true "Kredensial Login"
// @Success 200 {object} helper.APIResponse{data=model.AuthResponse}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Router /auth/login [post]
func docLogin() {}

// docGetProfile godoc
// @Summary Get Current Profile
// @Description Mendapatkan data profil pengguna yang sedang terotentikasi
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helper.APIResponse{data=model.UserResponse}
// @Failure 401 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /auth/me [get]
func docGetProfile() {}

// ============================================================================
// 3. CLINICS
// ============================================================================

// docListClinics godoc
// @Summary List Public Clinics
// @Description Mendapatkan daftar seluruh klinik hewan aktif yang terdaftar di PETRIVA
// @Tags Clinics
// @Produce json
// @Success 200 {object} helper.APIResponse{data=[]model.Clinic}
// @Failure 500 {object} helper.APIResponse
// @Router /clinics [get]
func docListClinics() {}

// docGetClinic godoc
// @Summary Get Clinic Detail
// @Description Mendapatkan detail informasi klinik hewan berdasarkan ID
// @Tags Clinics
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.Clinic}
// @Failure 404 {object} helper.APIResponse
// @Router /clinics/{id} [get]
func docGetClinic() {}

// docCreateClinic godoc
// @Summary Create Clinic (Admin Only)
// @Description Menambahkan data klinik hewan baru ke dalam sistem
// @Tags Clinics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateClinicInput true "Data Klinik Baru"
// @Success 201 {object} helper.APIResponse{data=model.Clinic}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /clinics [post]
func docCreateClinic() {}

// docUpdateClinic godoc
// @Summary Update Clinic (Admin Only)
// @Description Memperbarui informasi data klinik hewan yang sudah ada
// @Tags Clinics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Param request body model.UpdateClinicInput true "Data Pembaruan Klinik"
// @Success 200 {object} helper.APIResponse{data=model.Clinic}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /clinics/{id} [put]
func docUpdateClinic() {}

// docDeleteClinic godoc
// @Summary Deactivate Clinic (Admin Only)
// @Description Menonaktifkan klinik hewan dari sistem
// @Tags Clinics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic ID (UUID)"
// @Success 200 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /clinics/{id} [delete]
func docDeleteClinic() {}

// ============================================================================
// 4. SERVICES
// ============================================================================

// docListServices godoc
// @Summary List Services
// @Description Mendapatkan daftar layanan klinik hewan dengan opsi filter berdasarkan klinik
// @Tags Services
// @Produce json
// @Param clinic_id query string false "Filter berdasarkan Clinic ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=[]model.Service}
// @Failure 500 {object} helper.APIResponse
// @Router /services [get]
func docListServices() {}

// docGetService godoc
// @Summary Get Service Detail
// @Description Mendapatkan detail informasi layanan berdasarkan ID
// @Tags Services
// @Produce json
// @Param id path string true "Service ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.Service}
// @Failure 404 {object} helper.APIResponse
// @Router /services/{id} [get]
func docGetService() {}

// docCreateService godoc
// @Summary Create Service (Admin Only)
// @Description Menambahkan layanan klinik baru
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateServiceInput true "Data Layanan Baru"
// @Success 201 {object} helper.APIResponse{data=model.Service}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /services [post]
func docCreateService() {}

// docUpdateService godoc
// @Summary Update Service (Admin Only)
// @Description Memperbarui informasi layanan klinik
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Service ID (UUID)"
// @Param request body model.UpdateServiceInput true "Data Pembaruan Layanan"
// @Success 200 {object} helper.APIResponse{data=model.Service}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /services/{id} [put]
func docUpdateService() {}

// docDeleteService godoc
// @Summary Deactivate Service (Admin Only)
// @Description Menonaktifkan layanan klinik
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param id path string true "Service ID (UUID)"
// @Success 200 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /services/{id} [delete]
func docDeleteService() {}

// ============================================================================
// 5. VETERINARIANS
// ============================================================================

// docListVeterinarians godoc
// @Summary List Veterinarians
// @Description Mendapatkan daftar seluruh dokter hewan aktif dengan opsi filter klinik
// @Tags Veterinarians
// @Produce json
// @Param clinic_id query string false "Filter berdasarkan Clinic ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=[]model.UserResponse}
// @Failure 500 {object} helper.APIResponse
// @Router /veterinarians [get]
func docListVeterinarians() {}

// docGetVeterinarian godoc
// @Summary Get Veterinarian Detail
// @Description Mendapatkan detail informasi dokter hewan berdasarkan ID
// @Tags Veterinarians
// @Produce json
// @Param id path string true "Veterinarian ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.UserResponse}
// @Failure 404 {object} helper.APIResponse
// @Router /veterinarians/{id} [get]
func docGetVeterinarian() {}

// ============================================================================
// 6. USERS
// ============================================================================

// docListUsers godoc
// @Summary List Users (Admin Only)
// @Description Mendapatkan daftar seluruh pengguna dengan filter opsional berdasarkan role
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param role query string false "Filter berdasarkan Role (ADMIN, VETERINARIAN, PET_OWNER)"
// @Success 200 {object} helper.APIResponse{data=[]model.UserResponse}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /users [get]
func docListUsers() {}

// docGetUser godoc
// @Summary Get User Detail (Admin Only)
// @Description Mendapatkan detail informasi data pengguna berdasarkan ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.UserResponse}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /users/{id} [get]
func docGetUser() {}

// docCreateUser godoc
// @Summary Create User (Admin Only)
// @Description Menambahkan pengguna baru ke dalam sistem
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateUserInput true "Data Pengguna Baru"
// @Success 201 {object} helper.APIResponse{data=model.UserResponse}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /users [post]
func docCreateUser() {}

// docUpdateUser godoc
// @Summary Update User (Admin Only)
// @Description Memperbarui data profil atau status pengguna
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body model.UpdateUserInput true "Data Pembaruan Pengguna"
// @Success 200 {object} helper.APIResponse{data=model.UserResponse}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /users/{id} [put]
func docUpdateUser() {}

// docDeleteUser godoc
// @Summary Deactivate User (Admin Only)
// @Description Menonaktifkan akun pengguna
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /users/{id} [delete]
func docDeleteUser() {}

// ============================================================================
// 7. PETS
// ============================================================================

// docCreatePet godoc
// @Summary Register Pet (Pet Owner Only)
// @Description Mendaftarkan hewan peliharaan baru untuk akun pemilik yang sedang login
// @Tags Pets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreatePetInput true "Data Hewan Baru"
// @Success 201 {object} helper.APIResponse{data=model.Pet}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /pets [post]
func docCreatePet() {}

// docListPets godoc
// @Summary List Pets
// @Description Mendapatkan daftar hewan peliharaan (milik sendiri bagi Pet Owner, atau seluruh hewan bagi Admin)
// @Tags Pets
// @Security BearerAuth
// @Produce json
// @Param owner_id query string false "Filter berdasarkan Owner ID (Khusus Admin)"
// @Success 200 {object} helper.APIResponse{data=[]model.Pet}
// @Failure 401 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /pets [get]
func docListPets() {}

// docGetPet godoc
// @Summary Get Pet Detail
// @Description Mendapatkan informasi lengkap profil hewan peliharaan berdasarkan ID
// @Tags Pets
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pet ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.Pet}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /pets/{id} [get]
func docGetPet() {}

// docUpdatePet godoc
// @Summary Update Pet (Pet Owner Only)
// @Description Memperbarui data hewan peliharaan milik sendiri
// @Tags Pets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Pet ID (UUID)"
// @Param request body model.UpdatePetInput true "Data Pembaruan Hewan"
// @Success 200 {object} helper.APIResponse{data=model.Pet}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /pets/{id} [put]
func docUpdatePet() {}

// docDeletePet godoc
// @Summary Deactivate Pet (Pet Owner Only)
// @Description Menonaktifkan data hewan peliharaan milik sendiri
// @Tags Pets
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pet ID (UUID)"
// @Success 200 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /pets/{id} [delete]
func docDeletePet() {}

// ============================================================================
// 8. APPOINTMENTS
// ============================================================================

// docCreateAppointment godoc
// @Summary Book Appointment (Pet Owner Only)
// @Description Membuat jadwal janji temu pemeriksaan hewan ke klinik dan dokter
// @Tags Appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateAppointmentInput true "Data Janji Temu Baru"
// @Success 201 {object} helper.APIResponse{data=model.Appointment}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /appointments [post]
func docCreateAppointment() {}

// docListAppointments godoc
// @Summary List Appointments
// @Description Mendapatkan daftar janji temu sesuai role pengguna dengan opsi filter status
// @Tags Appointments
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter Status (PENDING, CONFIRMED, IN_PROGRESS, COMPLETED, CANCELLED, REJECTED)"
// @Success 200 {object} helper.APIResponse{data=[]model.Appointment}
// @Failure 401 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /appointments [get]
func docListAppointments() {}

// docGetAppointment godoc
// @Summary Get Appointment Detail
// @Description Mendapatkan detail janji temu berdasarkan ID
// @Tags Appointments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.Appointment}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /appointments/{id} [get]
func docGetAppointment() {}

// docCancelAppointment godoc
// @Summary Cancel Appointment
// @Description Membatalkan janji temu dengan menyertakan alasan pembatalan
// @Tags Appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body model.CancelAppointmentInput true "Alasan Pembatalan"
// @Success 200 {object} helper.APIResponse{data=model.Appointment}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /appointments/{id}/cancel [put]
func docCancelAppointment() {}

// docUpdateAppointmentStatus godoc
// @Summary Update Appointment Status (Veterinarian & Admin)
// @Description Mengubah status janji temu (CONFIRMED, IN_PROGRESS, COMPLETED, REJECTED)
// @Tags Appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body model.UpdateAppointmentStatusInput true "Status Baru"
// @Success 200 {object} helper.APIResponse{data=model.Appointment}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /appointments/{id}/status [put]
func docUpdateAppointmentStatus() {}

// ============================================================================
// 9. MEDICAL RECORDS
// ============================================================================

// docCreateMedicalRecord godoc
// @Summary Create Medical Record (Veterinarian Only)
// @Description Mencatat data rekam medis hasil pemeriksaan janji temu hewan
// @Tags Medical Records
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateMedicalRecordInput true "Data Rekam Medis"
// @Success 201 {object} helper.APIResponse{data=model.MedicalRecord}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /medical-records [post]
func docCreateMedicalRecord() {}

// docListMedicalRecords godoc
// @Summary List Medical Records
// @Description Mendapatkan daftar rekam medis (milik hewan sendiri bagi Pet Owner, atau seluruhnya bagi Dokter & Admin)
// @Tags Medical Records
// @Security BearerAuth
// @Produce json
// @Param pet_id query string false "Filter berdasarkan Pet ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=[]model.MedicalRecord}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /medical-records [get]
func docListMedicalRecords() {}

// docGetMedicalRecord godoc
// @Summary Get Medical Record Detail
// @Description Mendapatkan detail data rekam medis pemeriksaan berdasarkan ID
// @Tags Medical Records
// @Security BearerAuth
// @Produce json
// @Param id path string true "Medical Record ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.MedicalRecord}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /medical-records/{id} [get]
func docGetMedicalRecord() {}

// docUpdateMedicalRecord godoc
// @Summary Update Medical Record (Veterinarian & Admin)
// @Description Memperbarui diagnosa, terapi, resep, atau catatan rekam medis
// @Tags Medical Records
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Medical Record ID (UUID)"
// @Param request body model.UpdateMedicalRecordInput true "Data Pembaruan Rekam Medis"
// @Success 200 {object} helper.APIResponse{data=model.MedicalRecord}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /medical-records/{id} [put]
func docUpdateMedicalRecord() {}

// ============================================================================
// 10. VACCINATIONS
// ============================================================================

// docCreateVaccination godoc
// @Summary Record Vaccination (Veterinarian Only)
// @Description Mencatat data pemberian vaksin hewan peliharaan
// @Tags Vaccinations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateVaccinationInput true "Data Vaksinasi"
// @Success 201 {object} helper.APIResponse{data=model.Vaccination}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Router /vaccinations [post]
func docCreateVaccination() {}

// docListVaccinations godoc
// @Summary List Vaccinations
// @Description Mendapatkan riwayat daftar vaksinasi hewan peliharaan
// @Tags Vaccinations
// @Security BearerAuth
// @Produce json
// @Param pet_id query string false "Filter berdasarkan Pet ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=[]model.Vaccination}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /vaccinations [get]
func docListVaccinations() {}

// docGetVaccination godoc
// @Summary Get Vaccination Detail
// @Description Mendapatkan detail sertifikat/catatan vaksinasi berdasarkan ID
// @Tags Vaccinations
// @Security BearerAuth
// @Produce json
// @Param id path string true "Vaccination ID (UUID)"
// @Success 200 {object} helper.APIResponse{data=model.Vaccination}
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /vaccinations/{id} [get]
func docGetVaccination() {}

// docUpdateVaccination godoc
// @Summary Update Vaccination (Veterinarian & Admin)
// @Description Memperbarui data vaksinasi hewan
// @Tags Vaccinations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vaccination ID (UUID)"
// @Param request body model.UpdateVaccinationInput true "Data Pembaruan Vaksinasi"
// @Success 200 {object} helper.APIResponse{data=model.Vaccination}
// @Failure 400 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 404 {object} helper.APIResponse
// @Router /vaccinations/{id} [put]
func docUpdateVaccination() {}

// docDeleteVaccination godoc
// @Summary Delete Vaccination (Veterinarian & Admin)
// @Description Menghapus catatan riwayat vaksinasi
// @Tags Vaccinations
// @Security BearerAuth
// @Produce json
// @Param id path string true "Vaccination ID (UUID)"
// @Success 200 {object} helper.APIResponse
// @Failure 401 {object} helper.APIResponse
// @Failure 403 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /vaccinations/{id} [delete]
func docDeleteVaccination() {}

// ============================================================================
// 11. DASHBOARD
// ============================================================================

// docGetDashboardSummary godoc
// @Summary Get Dashboard Summary
// @Description Mendapatkan ringkasan statistik dan metrik operasional sesuai role pengguna (Admin, Dokter, atau Pet Owner)
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helper.APIResponse{data=model.DashboardSummary}
// @Failure 401 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /dashboard/summary [get]
func docGetDashboardSummary() {}
