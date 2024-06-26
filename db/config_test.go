package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRolesCheck_UserHasRole(t *testing.T) {
	// Mock user roles
	userRoles := []WorkspaceUserRoles{
		{Role: "ADD BOUNTY", OwnerPubKey: "user1", WorkspaceUuid: "org1", Created: &time.Time{}},
	}

	// Role to check
	roleToCheck := "ADD BOUNTY"

	// Call the function
	result := RolesCheck(userRoles, roleToCheck)

	// Assert that it returns true
	if !result {
		t.Errorf("Expected RolesCheck to return true for user with role %s, got false", roleToCheck)
	}
}

func TestRolesCheck_UserDoesNotHaveRole(t *testing.T) {
	// Mock user roles
	userRoles := []WorkspaceUserRoles{
		{Role: "DELETE BOUNTY", OwnerPubKey: "user2", WorkspaceUuid: "org1", Created: &time.Time{}},
	}

	// Role to check
	roleToCheck := "ADD BOUNTY"

	// Call the function
	result := RolesCheck(userRoles, roleToCheck)

	// Assert that it returns false
	if result {
		t.Errorf("Expected RolesCheck to return false for user without role %s, got true", roleToCheck)
	}
}

func TestCheckUser(t *testing.T) {
	userRoles := []WorkspaceUserRoles{
		{OwnerPubKey: "userPublicKey"},
	}

	// if in the user roles, one of the owner_pubkey belongs to the user return true else return false
	assert.True(t, CheckUser(userRoles, "userPublicKey"))
	assert.False(t, CheckUser(userRoles, "anotherPublicKey"))
}

func TestUserHasAccess(t *testing.T) {
	mockGetWorkspaceByUuid := func(uuid string) Workspace {
		return Workspace{
			Uuid:        uuid,
			OwnerPubKey: "org_admin",
		}
	}

	mockGetUserRoles := func(uuid string, pubkey string) []WorkspaceUserRoles {
		return []WorkspaceUserRoles{
			{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
		}
	}

	mockDB := &gorm.DB{}

	databaseConfig := NewDatabaseConfig(mockDB)
	databaseConfig.getWorkspaceByUuid = mockGetWorkspaceByUuid
	databaseConfig.getUserRoles = mockGetUserRoles

	t.Run("Should test that if the user is the admin of an workspace returns true", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("org_admin", "workspace_uuid", "ADD BOUNTY")

		// Assert that it returns true since the user is the org admin
		if !result {
			t.Errorf("Expected UserHasAccess to return true for workspace admin, got false")
		}
	})

	t.Run("Should test that if the user is not the workspace admin, and the user has the required role it should return true", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("user_pubkey", "workspace_uuid", "ADD BOUNTY")

		// Assert that it returns true since the user has the required role
		if !result {
			t.Errorf("Expected UserHasAccess to return true for user with required role, got false")
		}
	})

	t.Run("Should test that if the user is not the workspace admin, and the user has not the required role it should return false", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("user_pubkey", "workspace_uuid", "DELETE BOUNTY")

		// Assert that it returns false since the user does not have the required role
		if result {
			t.Errorf("Expected UserHasAccess to return false for user without required role, got true")
		}
	})
}

func TestUserHasManageBountyRoles(t *testing.T) {
	mockGetWorkspaceByUuid := func(uuid string) Workspace {
		return Workspace{
			Uuid:        uuid,
			OwnerPubKey: "org_admin",
		}
	}

	mockGetUserRoles := func(uuid string, pubkey string) []WorkspaceUserRoles {
		if uuid == "workspace_uuid" {
			return []WorkspaceUserRoles{
				{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
			}
		} else {
			return []WorkspaceUserRoles{
				{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "UPDATE BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "DELETE BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "PAY BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
			}
		}
	}

	mockDB := &gorm.DB{}

	databaseConfig := NewDatabaseConfig(mockDB)
	databaseConfig.getWorkspaceByUuid = mockGetWorkspaceByUuid
	databaseConfig.getUserRoles = mockGetUserRoles

	t.Run("Should test that if the user is the workspace admin return true", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("org_admin", "workspace_uuid")

		// Assert that it returns true since the user is the org admin
		assert.True(t, result, "Expected UserHasManageBountyRoles to return true for workspace admin")
	})

	t.Run("Should test that if the user has all bounty roles return true", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("user_pubkey", "workspace_uuid2")

		// Assert that it returns true since the user has all bounty roles
		assert.True(t, result, "Expected UserHasManageBountyRoles to return true for user with all bounty roles")
	})

	t.Run("Should test that if the user don't have all bounty roles return false.", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("user_pubkey", "workspace_uuid")

		// Assert that it returns false since the user does not have all bounty roles
		assert.False(t, result, "Expected UserHasManageBountyRoles to return false for user without all bounty roles")
	})
}
