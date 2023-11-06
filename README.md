# simple validator for map value

## example :

### Example 1
```go
payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
    "jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
    "hoby":          {Type: reflect.String, Null: false},
    "menikah":       {Type: reflect.Bool, Null: false},
}).Load(payload).RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}

err = map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
    "jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
    "hoby":          {Type: reflect.Int, Null: false},
    "menikah":       {Type: reflect.Bool, Null: false},
}).Load(payload).RunValidate()
if err == nil {
    t.Error("Expected have an error, but you got no error")
}
```

### Example Echo Framework
```go
func handleLogin(c echo.Context) error {
    jsonHttp, err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
        "email":    {Email: true, Max: map_validator.ToPointer[int](100)},
        "password": {Type: reflect.String, Min: map_validator.ToPointer[int](6), Max: map_validator.ToPointer[int](30)},
    }).LoadJsonHttp(c.Request())
    if err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    err = jsonHttp.RunValidate()
    if err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    return c.NoContent(http.StatusOK)
}

func main() {
    e := echo.New()
    e.POST("/login", handleLogin)
    e.Start(":3000")
}

```