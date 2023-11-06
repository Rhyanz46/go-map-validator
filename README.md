# simple validator for map value

## example :
### Gin Framework
```go
type MessageResponse struct {
    Message string `json:"message"`
}

func Create(c *gin.Context){
    var data map[string]interface{}
    err := json.NewDecoder(c.Request.Body).Decode(&data)
    if err != nil {
        if err.Error() == "EOF" {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: "please input a json data"})
            return
        }
        c.JSON(http.StatusBadRequest, MessageResponse{Message: "json format is invalid"})
        return
    }
    statusCode, err := ValidateJson(data)
    if err != nil{
        c.JSON(statusCode, MessageResponse{Message: err.Error()})
        return
    }
    c.JSON(http.StatusBadRequest, MessageResponse{Message: "ok"})
}

func ValidateJson(payload map[string]interface{}) (statusCode int, err error) {
    _, err = mapValidator.Validate(
        "project_id", payload, mapValidator.RequestDataValidator{UUID: true, Null: false},
    )
    if err != nil {
        return http.StatusBadRequest, err
    }
    _, err = mapValidator.Validate(
        "name", payload, mapValidator.RequestDataValidator{Type: reflect.String, Null: false, Max: mapValidator.ToPointer[int](100)},
    )
    if err != nil {
        return http.StatusBadRequest, err
    }
    _, err := mapValidator.Validate(
        "description", payload, mapValidator.RequestDataValidator{Type: reflect.String, Null: true, NilIfNull: true},
    )
    if err != nil {
        return http.StatusBadRequest, err
    }
    return http.StatusOK, nil
}
```

### Enum Validation
```go
payload := map[string]interface{}{"data": "arian", "jenis_kelamin": "laki-laki", "hoby": "Main PS"}
_, err := Validate(
    "data", payload, RequestDataValidator{
        Null: false,
        Enum: &EnumField[any]{Items: []string{"arian", "aaa"}},
    },
)
if err != nil {
    t.Errorf("Test case 1 Error : %v", err)
}
```

### Email Validation
```go
payload := map[string]interface{}{
    "email": "test@example.com",
}

validator := RequestDataValidator{
    Email: true,
}

_, err := Validate("email", payload, validator)

if err != nil {
    fmt.Println("Error :", err)
}
```


### IPv4 Validation
```go
payload := map[string]interface{}{
    "ip_address": "192.168.1.1",
}

validator := RequestDataValidator{
    IPV4: true,
}

_, err := Validate("email", payload, validator)

if err != nil {
    fmt.Println("Error :", err)
}
```

### Multiple Error Handling
```go
payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
err := MultiValidate(payload, map[string]RequestDataValidator{
    "jenis_kelamin": {Enum: &EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
    "hoby":          {Type: reflect.String, Null: false},
    "menikah":       {Type: reflect.Bool, Null: false},
})
if err != nil {
    fmt.Println("Expected not have error, but got error : ", err)
}
```
