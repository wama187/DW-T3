package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
    secretKey []byte
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
    return &AuthMiddleware{
        secretKey: []byte(secretKey),
    }
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Obtener el token del header Authorization
        authHeader := r.Header.Get("Authorization")
        
        if authHeader != "" {
            // Formato esperado: "Bearer <token>"
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                tokenString := parts[1]
                
                // Parsear y validar el token
                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                    // Verificar que el método de firma sea el esperado
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, jwt.ErrSignatureInvalid
                    }
                    return m.secretKey, nil
                })
                
                if err == nil && token.Valid {
                    // Extraer los claims
                    if claims, ok := token.Claims.(jwt.MapClaims); ok {
                        // Extraer el user_id
                        if userID, ok := claims["user_id"].(float64); ok {
                            // Añadir el userID al contexto
                            ctx := context.WithValue(r.Context(), UserContextKey, int(userID))
                            r = r.WithContext(ctx)
                        }
                    }
                }
            }
        }
        
        next.ServeHTTP(w, r)
    })
}

func GetUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(UserContextKey).(int)
    return userID, ok
}