// Package middleware contains HTTP middleware implementations.
//
// This file implements role-based access control (RBAC) middleware.
// It performs permission-based authorization checks at the route level.
// Middleware extracts user permissions and validates against required
// permissions for each protected endpoint.
package middleware