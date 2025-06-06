package services

import (
	"libreria/models"
	"libreria/repositories"
	"libreria/utils"
	"log"
	"os"
)

func SeedInitialData(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) error {
	// Seed roles
	defaultRoles := []string{"ADMIN", "WRITE", "READ", "ROOT"}

	for _, roleName := range defaultRoles {
		_, err := roleRepo.FindByRoleName(roleName)
		if err != nil {
			if err := roleRepo.Create(&models.Role{Name: roleName}); err != nil {
				return err
			}
		}
	}

	// Seed ROOT user
	username := "root"
	password := os.Getenv("ROOT_PASSWORD") // asegurate de cargarlo antes con godotenv
	email := os.Getenv("ROOT_EMAIL")

	if email == "" || password == "" {
		log.Println("ROOT_EMAIL o ROOT_PASSWORD no seteados, se saltea creación de usuario root")
		return nil
	}

	_, err := userRepo.FindByUserName(username)
	if err == nil {
		log.Println("Usuario root ya existe, se omite creación.")
		return nil
	}

	hashedPass, _ := utils.HashPassword(password)

	rootUser := models.User{
		Username: username,
		Email:    email,
		Password: hashedPass,
	}

	if err := userRepo.Create(&rootUser); err != nil {
		return err
	}

	// Asignar rol ROOT
	adminRole, _ := roleRepo.FindByRoleName("ROOT")
	userRole := models.UserRole{
		UserID: rootUser.ID,
		RoleID: adminRole.ID,
	}

	return userRepo.CreateUserRole(&userRole)

}
