package users

import (
    "context"
    "fmt"
    "strconv"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    
    "task-manager/ent"
    "task-manager/ent/user"
    "task-manager/graph/model"
)

type UserService struct {
    client    *ent.Client
    secretKey []byte
}

func NewUserService(client *ent.Client, secretKey string) *UserService {
    return &UserService{
        client:    client,
        secretKey: []byte(secretKey),
    }
}

func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*model.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("error al hashear contraseña: %w", err)
    }
    
    entUser, err := s.client.User.
        Create().
        SetUsername(username).
        SetEmail(email).
        SetPassword(string(hashedPassword)).
        Save(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al crear usuario: %w", err)
    }
    
    return &model.User{
        ID:       strconv.Itoa(entUser.ID),
        Username: entUser.Username,
        Email:    entUser.Email,
    }, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*model.AuthResponse, error) {
    entUser, err := s.client.User.
        Query().
        Where(user.UsernameEQ(username)).
        Only(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("credenciales inválidas")
    }
    
    // Verificar contraseña
    err = bcrypt.CompareHashAndPassword([]byte(entUser.Password), []byte(password))
    if err != nil {
        return nil, fmt.Errorf("credenciales inválidas")
    }
    
    // Generar JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": entUser.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    
    tokenString, err := token.SignedString(s.secretKey)
    if err != nil {
        return nil, fmt.Errorf("error al generar token: %w", err)
    }
    
    return &model.AuthResponse{
        Token: tokenString,
        User: &model.User{
            ID:       strconv.Itoa(entUser.ID),
            Username: entUser.Username,
            Email:    entUser.Email,
        },
    }, nil
}