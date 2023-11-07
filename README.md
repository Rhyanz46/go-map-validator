# simple validator for map value

You can integrate with all golang framework using `http.Request` interface

examples : https://github.com/Rhyanz46/go-map-validator-examples

### Install 
```shell
go get github.com/Rhyanz46/go-map-validator/map_validator
```

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

### Example Bind To Struct
```go
type Data struct {
    JK      string `map_validator:"jenis_kelamin"`
    Hoby    string `map_validator:"hoby"`
    Menikah bool   `map_validator:"menikah"`
}

payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
    "jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
    "hoby":          {Type: reflect.String, Null: false},
    "menikah":       {Type: reflect.Bool, Null: false},
}).Load(payload).RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}

testBind := &Data{}
if testBind.JK != "" {
    t.Errorf("Expected : '' But you got : %s", testBind.JK)
}
err = extraCheck.Bind(testBind)
if err != nil {
    t.Errorf("Error : %s ", err)
}

if testBind.JK != payload["jenis_kelamin"] {
    t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
}

```

## Features
- validate value in by keys `map[string]interface{}`
- validate data from `http.Request` json/multipart
  - support file upload

## Road Map
- avoiding same value in some field
- get from urls params http
- validation for `base64`
- handle file size on multipart
- generate OpenAPI Spec