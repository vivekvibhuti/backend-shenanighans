// Package middleware contains HTTP middleware implementations.
//
// This file implements authentication middleware that extracts and validates
// user identity from requests. It populates the request context with user
// information for downstream handlers and middleware.
package middleware