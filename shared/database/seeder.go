package database

import (
	"database/sql"
	"errors"
	"log"
	"os/user"
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/shared/config"
	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) {
	tx := db.Begin()

	{ // Seeding Module
		var totalModule int64
		tx.Model(&entity.Module{}).Where("name IN ?", []string{"User", "Role", "Permission", "Module"}).Count(&totalModule)
		if totalModule != 4 {
			if err := seedingModuleUserManagement(tx); err != nil {
				log.Printf("Seeding module user management failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding module user management")
		}
	}

	{ // Seeding permission user management
		var totalPermission int64
		tx.Model(&entity.Permission{}).Count(&totalPermission)
		if totalPermission == 0 {
			if err := seedingPermissionUserManagement(tx); err != nil {
				log.Printf("Seeding permission user management failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding permission user management")
		}
	}

	{ // Seeding role admin
		var totalRoleAdmin int64
		tx.Model(&entity.Role{}).Where("name = ?", "Admin").Count(&totalRoleAdmin)
		if totalRoleAdmin == 0 {
			if err := seedingRoleAdmin(tx); err != nil {
				log.Printf("Seeding role admin failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding role admin")
		}
	}

	{ // Seeding role guest
		var totalRoleUser int64
		tx.Model(&entity.Role{}).Where("name = ?", "user").Count(&totalRoleUser)
		if totalRoleUser == 0 {
			if err := seedingRoleGuest(tx); err != nil {
				log.Printf("Seeding role user failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding role user")
		}
	}

	{ // Seeding user admin
		var totalAdmin int64
		tx.Model(&user.User{}).Where("username = ? AND email = ?", "admin", "admin@email.id").Count(&totalAdmin)
		if totalAdmin == 0 {
			if err := seedingUserAdmin(tx); err != nil {
				log.Printf("Seeding user admin failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding user admin")
		}
	}

	{ // Seeding user guest
		var totalUser int64
		tx.Model(&user.User{}).Where("username = ? AND email = ?", "user", "user@email.id").Count(&totalUser)
		if totalUser == 0 {
			if err := seedingUserGuest(tx); err != nil {
				log.Printf("Seeding user guest failed: %v", err)
				tx.Rollback()
				return
			}

			log.Println("Success seeding user guest")
		}
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v", err)
	}
}

func seedingModuleUserManagement(tx *gorm.DB) error {
	modules := []entity.Module{}

	var userModule int64
	tx.Model(&entity.Module{}).Where("name = ?", "User").Count(&userModule)
	if userModule == 0 {
		userUUID, err := uuid.Parse("1234f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return err
		}

		modules = append(modules, entity.Module{
			Name: "User",
			UUID: userUUID,
		})

		log.Println("Seeding Module User")
	}

	var roleModule int64
	tx.Model(&entity.Module{}).Where("name = ?", "Role").Count(&roleModule)
	if roleModule == 0 {
		roleUUID, err := uuid.Parse("1235f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return err
		}

		modules = append(modules, entity.Module{
			Name: "Role",
			UUID: roleUUID,
		})

		log.Println("Seeding Module Role")
	}

	var permissionModule int64
	tx.Model(&entity.Module{}).Where("name = ?", "Permission").Count(&permissionModule)
	if permissionModule == 0 {
		permissionUUID, err := uuid.Parse("1236f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return err
		}

		modules = append(modules, entity.Module{
			Name: "Permission",
			UUID: permissionUUID,
		})

		log.Println("Seeding Module Permission")
	}

	var moduleModule int64
	tx.Model(&entity.Module{}).Where("name = ?", "Module").Count(&moduleModule)
	if moduleModule == 0 {
		moduleUUID, err := uuid.Parse("1237f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return err
		}

		modules = append(modules, entity.Module{
			Name: "Module",
			UUID: moduleUUID,
		})

		log.Println("Seeding Module Module")
	}

	if err := tx.Create(&modules).Error; err != nil {
		return err
	}

	return nil
}

func seedingPermissionUserManagement(tx *gorm.DB) error {
	permissions := []entity.Permission{}

	// Permission for user module
	{
		var userModule entity.Module
		if result := tx.Limit(1).Where("name = ?", "User").Find(&userModule); result.RowsAffected == 0 {
			return errors.New("user module not found")
		}

		userViewUUID, err := uuid.Parse("1238f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		userCreateUUID, err := uuid.Parse("1239f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		userUpdateUUID, err := uuid.Parse("1240f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		userDeleteUUID, err := uuid.Parse("1241f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}

		userPermissions := []entity.Permission{
			{
				UUID:     userViewUUID,
				Name:     "View User",
				ModuleID: userModule.ID,
			},
			{
				UUID:     userCreateUUID,
				Name:     "Create User",
				ModuleID: userModule.ID,
			},
			{
				UUID:     userUpdateUUID,
				Name:     "Update User",
				ModuleID: userModule.ID,
			},
			{
				UUID:     userDeleteUUID,
				Name:     "Delete User",
				ModuleID: userModule.ID,
			},
		}

		permissions = append(permissions, userPermissions...)
	}

	// Permission for role module
	{
		var roleModule entity.Module
		if result := tx.Limit(1).Where("name = ?", "Role").Find(&roleModule); result.RowsAffected == 0 {
			return errors.New("role module not found")
		}

		roleViewUUID, err := uuid.Parse("1242f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		roleCreateUUID, err := uuid.Parse("1243f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		roleUpdateUUID, err := uuid.Parse("1244f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		roleDeleteUUID, err := uuid.Parse("1245f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}

		rolePermissions := []entity.Permission{
			{
				UUID:     roleViewUUID,
				Name:     "View Role",
				ModuleID: roleModule.ID,
			},
			{
				UUID:     roleCreateUUID,
				Name:     "Create Role",
				ModuleID: roleModule.ID,
			},
			{
				UUID:     roleUpdateUUID,
				Name:     "Update Role",
				ModuleID: roleModule.ID,
			},
			{
				UUID:     roleDeleteUUID,
				Name:     "Delete Role",
				ModuleID: roleModule.ID,
			},
		}

		permissions = append(permissions, rolePermissions...)
	}

	// Permission for permission module
	{
		var permissionModule entity.Module
		if result := tx.Limit(1).Where("name = ?", "Permission").Find(&permissionModule); result.RowsAffected == 0 {
			return errors.New("permission module not found")
		}

		permissionViewUUID, err := uuid.Parse("1246f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		permissionCreateUUID, err := uuid.Parse("1247f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		permissionUpdateUUID, err := uuid.Parse("1248f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		permissionDeleteUUID, err := uuid.Parse("1249f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}

		permissionPermissions := []entity.Permission{
			{
				UUID:     permissionViewUUID,
				Name:     "View Permission",
				ModuleID: permissionModule.ID,
			},
			{
				UUID:     permissionCreateUUID,
				Name:     "Create Permission",
				ModuleID: permissionModule.ID,
			},
			{
				UUID:     permissionUpdateUUID,
				Name:     "Update Permission",
				ModuleID: permissionModule.ID,
			},
			{
				UUID:     permissionDeleteUUID,
				Name:     "Delete Permission",
				ModuleID: permissionModule.ID,
			},
		}

		permissions = append(permissions, permissionPermissions...)
	}

	// Permission for module module
	{
		var moduleModule entity.Module
		if result := tx.Limit(1).Where("name = ?", "Permission").Find(&moduleModule); result.RowsAffected == 0 {
			return errors.New("module module not found")
		}

		moduleViewUUID, err := uuid.Parse("1250f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		moduleCreateUUID, err := uuid.Parse("1251f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		moduleUpdateUUID, err := uuid.Parse("1252f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}
		moduleDeleteUUID, err := uuid.Parse("1253f6bf-8a3d-46de-a89d-ed901f90a7ad")
		if err != nil {
			return nil
		}

		modulePermissions := []entity.Permission{
			{
				UUID:     moduleViewUUID,
				Name:     "View Permission",
				ModuleID: moduleModule.ID,
			},
			{
				UUID:     moduleCreateUUID,
				Name:     "Create Permission",
				ModuleID: moduleModule.ID,
			},
			{
				UUID:     moduleUpdateUUID,
				Name:     "Update Permission",
				ModuleID: moduleModule.ID,
			},
			{
				UUID:     moduleDeleteUUID,
				Name:     "Delete Permission",
				ModuleID: moduleModule.ID,
			},
		}

		permissions = append(permissions, modulePermissions...)
	}

	if err := tx.Create(&permissions).Error; err != nil {
		return err
	}

	return nil
}

func seedingRoleAdmin(tx *gorm.DB) error {
	adminUUID, err := uuid.Parse("1254f6bf-8a3d-46de-a89d-ed901f90a7ad")
	if err != nil {
		return err
	}

	adminRole := entity.Role{
		UUID:    adminUUID,
		Name:    "Admin",
		IsAdmin: true,
	}

	if err := tx.Create(&adminRole).Error; err != nil {
		return err
	}

	// Get all permission
	var permissions []entity.Permission
	tx.Model(&entity.Permission{}).Find(&permissions)

	// Append all permissions to many to many table "role_permissions"
	tx.Model(&adminRole).Association("Permissions").Append(&permissions)

	return nil
}

func seedingRoleGuest(tx *gorm.DB) error {
	userUUID, err := uuid.Parse("4ff46fec-78ec-4f68-8db8-a495fac37c03")
	if err != nil {
		return err
	}

	userRole := entity.Role{
		UUID:    userUUID,
		Name:    "User",
		IsAdmin: false,
	}

	if err := tx.Create(&userRole).Error; err != nil {
		return err
	}

	// Get permission that user could use
	var permissions []entity.Permission
	tx.Model(&entity.Permission{}).Where("name IN ?", []string{"View User", "Update User"}).Find(&permissions)

	// Append related permissions to many to many table "role_permissions"
	tx.Model(&userRole).Association("Permissions").Append(&permissions)

	return nil
}

func seedingUserAdmin(tx *gorm.DB) error {
	cfg := config.AppConfig

	// Set few value for user admin
	adminUUID, err := uuid.Parse("3685f6bf-8a3d-46de-a89d-ed901f90a7ad")
	if err != nil {
		return err
	}

	// Find admin role
	var adminRole entity.Role
	if result := tx.Limit(1).Where("name = ?", "Admin").Find(&adminRole); result.RowsAffected == 0 {
		return errors.New("admin role not found")
	}

	user := entity.User{
		UUID:        adminUUID,
		RoleID:      adminRole.ID,
		Username:    "admin",
		Email:       "admin@email.id",
		Password:    cfg.AdminPass,
		ValidatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := tx.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func seedingUserGuest(tx *gorm.DB) error {
	// Set few value for user guest
	userUUID, err := uuid.Parse("b2db4155-a1e4-42d7-b5b5-415bcfe54cdd")
	if err != nil {
		return err
	}

	// Find guest role
	var userRole entity.Role
	if result := tx.Limit(1).Where("name = ?", "User").Find(&userRole); result.RowsAffected == 0 {
		return errors.New("user role not found")
	}

	user := entity.User{
		UUID:        userUUID,
		RoleID:      userRole.ID,
		Username:    "user",
		Email:       "user@email.id",
		Password:    "1234567",
		ValidatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := tx.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
