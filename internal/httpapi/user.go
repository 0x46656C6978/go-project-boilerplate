package httpapi

import (
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHttpApi struct {
	s service.UserServiceInterface
	cfg *config.Config
}

func NewUserHttpApi(cfg *config.Config, s service.UserServiceInterface) *UserHttpApi {
	return &UserHttpApi{
		s: s,
		cfg: cfg,
	}
}

func (u *UserHttpApi) RegisterApiEndpoints(engine *gin.Engine) {
	engine.POST("/users/login", u.Login)
	engine.POST("/users/register", u.Register)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserHttpApi) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}
	user, err := u.s.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "user not found",
		})
		return
	}
	err = u.s.VerifyCredentials(c, user, req.Email, req.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid credentials",
		})
		return
	}

	// generate jwt token
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(u.cfg.JWT.Expire)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    u.cfg.JWT.Issuer,
		Subject:   user.Email,
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := tok.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"token": signedStr,
		},
	})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserHttpApi) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}
	user, err := u.s.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}
	if user != nil {
		c.JSON(400, gin.H{
			"message": "user already exists",
		})
		return
	}

	user.Email = req.Email
	user.SetPassword(req.Password)
	err = u.s.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}
