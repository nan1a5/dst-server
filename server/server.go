package server

import (
	"strings"
	"time"

	"dst-manager/manager"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret = []byte("secret")
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

var users = map[string]User{
	"admin": {Username: "admin", Password: hash("admin123"), Role: "admin"},
}

func login(c *gin.Context) {
	var req struct {
		User string
		Pass string
	}
	c.BindJSON(&req)

	u, ok := users[req.User]
	if !ok || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Pass)) != nil {
		c.JSON(401, Response{
			Error:   "invalid",
			Status:  401,
			Message: "用户名或密码错误",
		})
		return
	}

	t, _ := tokenFor(u)
	c.JSON(200, Response{
		Data:    gin.H{"token": t},
		Status:  200,
		Message: "登录成功",
	})
}

func start_server(c *gin.Context) {
	manager.NewManager().Log("Starting server...")
	if err := manager.NewManager().StartServer(); err != nil {
		c.JSON(500, Response{
			Error:   "start_server_error",
			Status:  500,
			Message: "服务器启动失败: " + err.Error(),
		})
		return
	}
	c.JSON(200, Response{
		Status:  200,
		Message: "服务器启动成功",
	})
}

func server() *gin.Engine {
	r := gin.Default()
	r.POST("/login", login)

	api := r.Group("/api", auth())
	{
		api.POST("/start_server", start_server)
	}
	return r
}

// 工具函数
func tokenFor(user User) (string, error) {
	claims := jwt.MapClaims{
		"user": user.Username,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		token, err := jwt.Parse(tk, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}

func hash(p string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(b)
}
