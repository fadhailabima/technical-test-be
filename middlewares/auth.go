package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware - Middleware untuk memvalidasi JWT token pada setiap request
// Fungsi: Extract token dari header -> Validasi token -> Simpan user info ke context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Authorization header dari request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		// 2. Validasi format header harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Parse dan validasi JWT token dengan secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method harus HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// 4. Cek apakah token valid dan tidak expired
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			c.Abort()
			return
		}

		// 5. Extract claims (user info) dari token dan simpan ke gin context
		// Data ini bisa diakses oleh handler berikutnya
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["sub"])  // User ID dari claim "sub"
			c.Set("role", claims["role"])   // Role dari claim "role"
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid claims"})
			c.Abort()
			return
		}

		// 6. Lanjutkan ke handler berikutnya
		c.Next()
	}
}

// RoleMiddleware - Middleware untuk authorization berdasarkan role user
// Parameter: allowedRoles - daftar role yang diizinkan akses endpoint
// Contoh: RoleMiddleware("Admin", "Seller") -> hanya Admin dan Seller bisa akses
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil role yang sudah di-set oleh AuthMiddleware sebelumnya
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role user tidak ditemukan"})
			c.Abort()
			return
		}

		roleString := userRole.(string)
		isAllowed := false

		// 2. Cek apakah role user ada dalam daftar allowedRoles
		for _, role := range allowedRoles {
			if role == roleString {
				isAllowed = true
				break
			}
		}

		// 3. Jika role tidak diizinkan, return 403 Forbidden
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses (Forbidden)"})
			c.Abort()
			return
		}

		// 4. Role sesuai, lanjutkan ke handler
		c.Next()
	}
}