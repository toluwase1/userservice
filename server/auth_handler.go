package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-service/errors"
	"user-service/models"
	"user-service/server/jwt"
	"user-service/server/response"
	"user-service/services"
)

const AuthHeader = "Authorization"

func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get(AuthHeader)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := decode(c, &user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Signup successful", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := decode(c, &loginRequest); err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		userResponse, err := s.AuthService.LoginUser(&loginRequest)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		response.JSON(c, "login successful", http.StatusOK, userResponse, nil)
	}
}

func (s *Server) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := TokenFromHeader(c.Request)
		if accessToken == "" {
			response.Unauthorized(c, "access token required")
			return
		}

		token, err := jwt.ValidateAndGetClaims(accessToken, services.Secrete)
		if err != nil {
			response.Unauthorized(c, err.Error())
			return
		}

		if s.AuthService.IsTokenBlacklisted(accessToken) {
			response.Unauthorized(c, "invalid access token")
			return
		}

		response.JSON(c, "authenticated", http.StatusOK, token, nil)
		c.Next()
	}
}
