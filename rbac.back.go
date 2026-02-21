//go:build !wasm

package site

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/rbac"
	"github.com/tinywasm/unixid"
)

// Type aliases — callers don't need to import tinywasm/rbac directly.
type DBExecutor = rbac.Executor // Exec + QueryRow + Query
type DBScanner = rbac.Scanner   // Scan(dest ...any) error
type DBRows = rbac.Rows         // Next/Scan/Close/Err

type roleSpec struct {
	code        byte
	name        string
	description string
}

var (
	dbExecutor      DBExecutor
	getUserID       func(data ...any) string
	pendingRoles    []roleSpec
	pendingHandlers []any
	rbacInitialized bool
)

// SetDB sets the database executor. rbac initialization is deferred to Serve/Mount.
func SetDB(exec DBExecutor) {
	dbExecutor = exec
}

// SetUserID configures how to extract the current user's ID from request data.
// Required when SetDB has been called. Validated at Mount time.
func SetUserID(fn func(data ...any) string) {
	getUserID = fn
}

// CreateRole queues a role for creation at Serve/Mount time.
// Idempotent: safe to call on every startup (ON CONFLICT (code) DO NOTHING).
func CreateRole(code byte, name, description string) {
	pendingRoles = append(pendingRoles, roleSpec{code, name, description})
}

// AssignRole grants a role (identified by code) to a user.
// Typically called in the login handler after authentication.
func AssignRole(userID string, roleCode byte) error {
	if !rbacInitialized {
		return fmt.Err("site: Serve must be called before AssignRole")
	}
	role, err := rbac.GetRoleByCode(roleCode)
	if err != nil {
		return err
	}
	return rbac.AssignRole(userID, role.ID)
}

// RevokeRole removes a role (identified by code) from a user.
func RevokeRole(userID string, roleCode byte) error {
	if !rbacInitialized {
		return fmt.Err("site: Serve must be called before RevokeRole")
	}
	role, err := rbac.GetRoleByCode(roleCode)
	if err != nil {
		return err
	}
	return rbac.RevokeRole(userID, role.ID)
}

// GetUserRoleCodes returns the role codes assigned to a user (e.g., []byte{'a', 'e'}).
func GetUserRoleCodes(userID string) ([]byte, error) {
	if !rbacInitialized {
		return nil, fmt.Err("site: Serve must be called before GetUserRoleCodes")
	}
	return rbac.GetUserRoleCodes(userID)
}

// registerRBAC queues handlers for permission seeding. Applied by applyRBAC at Mount time.
func registerRBAC(handlers ...any) error {
	pendingHandlers = append(pendingHandlers, handlers...)
	return nil
}

// applyRBAC initializes rbac and seeds roles and permissions from queued state.
// Called once at Mount time. No-op when SetDB was not called (dev mode).
func applyRBAC() error {
	if dbExecutor == nil {
		return nil
	}
	if err := rbac.Init(dbExecutor); err != nil {
		return err
	}
	rbacInitialized = true
	handler.cp.SetAccessCheck(func(resource string, action byte, data ...any) bool {
		if getUserID == nil {
			return false
		}
		userID := getUserID(data...)
		if userID == "" {
			return false
		}
		ok, _ := rbac.HasPermission(userID, resource, action)
		return ok
	})
	for _, r := range pendingRoles {
		u, err := unixid.NewUnixID()
		if err != nil {
			return err
		}
		if err := rbac.CreateRole(u.GetNewID(), r.code, r.name, r.description); err != nil {
			return err
		}
	}
	if len(pendingHandlers) > 0 {
		return rbac.Register(pendingHandlers...)
	}
	return nil
}
