# Generic Mapper Usage

The Generic Mapper provides a flexible, type-safe way to convert between different struct types using Go 1.18+ generics and reflection.

## Basic Usage

### Simple Field Mapping
```go
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Create mapper and convert
mapper := NewFieldMapper[*User, *UserResponse]()
user := &User{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30}
response := mapper.Map(user)
// Result: UserResponse{ID: 1, Name: "John Doe", Email: "john@example.com"}
```

### Custom Field Mapping
```go
type Product struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    PriceCents  int64   `json:"price_cents"`
    CategoryID  int64   `json:"category_id"`
}

type ProductResponse struct {
    ID           int64  `json:"id"`
    Name         string `json:"name"`
    Price        string `json:"price"`        // Custom formatting
    CategoryName string `json:"category_name"` // Lookup required
}

// Create mapper with custom field logic
mapper := NewFieldMapper[*Product, *ProductResponse]().
    WithCustomMapping("Price", func(p *Product) any {
        return fmt.Sprintf("$%.2f", float64(p.PriceCents)/100)
    }).
    WithCustomMapping("CategoryName", func(p *Product) any {
        return categoryService.GetNameByID(p.CategoryID)
    })

product := &Product{ID: 1, Name: "Laptop", PriceCents: 99999, CategoryID: 5}
response := mapper.Map(product)
// Result: ProductResponse{ID: 1, Name: "Laptop", Price: "$999.99", CategoryName: "Electronics"}
```

## Advanced Examples

### Request to Model Mapping
```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type User struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Create request-to-model mapper
mapper := NewFieldMapper[*CreateUserRequest, *User]().
    WithCustomMapping("PasswordHash", func(req *CreateUserRequest) any {
        hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        return string(hash)
    }).
    WithCustomMapping("CreatedAt", func(req *CreateUserRequest) any {
        return time.Now()
    }).
    WithCustomMapping("UpdatedAt", func(req *CreateUserRequest) any {
        return time.Now()
    })
```

### Builder Mapper for Complex Logic
```go
// For cases requiring full control over mapping logic
orderMapper := NewBuilderMapper[*Order, *OrderResponse](func(order *Order) *OrderResponse {
    return &OrderResponse{
        ID:          order.ID,
        CustomerID:  order.CustomerID,
        TotalAmount: calculateTotal(order.Items),
        Status:      order.Status.String(),
        Items:       MapSlice(order.Items, itemToResponseMapper),
        CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
    }
})
```

### Slice Mapping Utility
```go
// Generic helper for slice conversions
func MapSlice[S any, T any](source []S, mapper Mapper[S, T]) []T {
    if source == nil {
        return nil
    }
    
    result := make([]T, len(source))
    for i, item := range source {
        result[i] = mapper.Map(item)
    }
    return result
}

// Usage
users := []*User{user1, user2, user3}
responses := MapSlice(users, userToResponseMapper)
```

## Creating Domain-Specific Mappers

### Complete Mapper Service
```go
type UserMapper struct {
    toResponse   Mapper[*User, *UserResponse]
    fromCreate   Mapper[*CreateUserRequest, *User]
    fromUpdate   Mapper[*UpdateUserRequest, *User]
}

func NewUserMapper() *UserMapper {
    return &UserMapper{
        toResponse: NewFieldMapper[*User, *UserResponse](),
        fromCreate: NewFieldMapper[*CreateUserRequest, *User]().
            WithCustomMapping("PasswordHash", hashPassword).
            WithCustomMapping("CreatedAt", setCurrentTime),
        fromUpdate: NewFieldMapper[*UpdateUserRequest, *User]().
            WithCustomMapping("UpdatedAt", setCurrentTime),
    }
}

// Convenience methods
func (um *UserMapper) ToResponse(user *User) *UserResponse {
    return um.toResponse.Map(user)
}

func (um *UserMapper) FromCreateRequest(req *CreateUserRequest) *User {
    return um.fromCreate.Map(req)
}
```

## Real-World Examples

### User Profile Mapper with Privacy Controls
```go
type ProfileMapper struct {
    toPublicProfile  Mapper[*User, *PublicProfile]
    toPrivateProfile Mapper[*User, *PrivateProfile]
    toAdminProfile   Mapper[*User, *AdminProfile]
}

func NewProfileMapper() *ProfileMapper {
    return &ProfileMapper{
        toPublicProfile: NewFieldMapper[*User, *PublicProfile]().
            WithCustomMapping("AvatarURL", func(u *User) any {
                if u.AvatarPath != "" {
                    return "/avatars/" + u.AvatarPath
                }
                return "/avatars/default.png"
            }).
            WithCustomMapping("JoinedDate", func(u *User) any {
                return u.CreatedAt.Format("January 2006")
            }),
            
        toPrivateProfile: NewFieldMapper[*User, *PrivateProfile]().
            WithCustomMapping("LastLoginFormatted", func(u *User) any {
                if u.LastLoginAt != nil {
                    return u.LastLoginAt.Format("2006-01-02 15:04:05")
                }
                return "Never"
            }),
            
        toAdminProfile: NewFieldMapper[*User, *AdminProfile](), // All fields mapped
    }
}
```

## Migration from Manual Mapping

### Before (Manual Mapping)
```go
func ToUserResponse(user *User) *UserResponse {
    if user == nil {
        return nil
    }
    
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
}
```

### After (Generic Mapper)
```go
var userToResponseMapper = NewFieldMapper[*User, *UserResponse]()

func ToUserResponse(user *User) *UserResponse {
    return userToResponseMapper.Map(user)
}
```

The Generic Mapper eliminates boilerplate code while providing type safety and extensibility for complex mapping scenarios.
