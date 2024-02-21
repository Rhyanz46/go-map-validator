# simple validator using struct properties as rules

You can integrate with all golang framework using `http.Request` interface

examples : https://github.com/Rhyanz46/go-map-validator/tree/main/test

### Install

```shell
go get github.com/Rhyanz46/go-map-validator/map_validator
```

### discussion or updates
- [On Telegram](https://t.me/addlist/Wi84VFNkvz85MWFl)

## Features

- validate value in `map[string]interface{}` by keys
- validate data from `http.Request` json/multipart
    - support file upload
- enum value check
- min/max length data check
- email field check
- uuid field check
- IPv4 field check
- IPv4 Network check
- regex on string validation
- nested validation 🔥
- you can create your own extension 🔥🔥🔥🔥 (example : [https://github.com/Rhyanz46/go-map-validator/example_extensions/](https://github.com/Rhyanz46/go-map-validator/tree/main/example_extensions))
- custom message :
    - on invalid regex message : ✅
    - on type not match message : ✅
    - on null data message : ❌
    - on max data message : ❌
    - on enum value not match : ❌

## On Progress

- validation for one data value only

## Custom Message Variables

| No |   Variable Name    |
|:--:|:------------------:| 
| 1  |     `${field}`     |
| 2  | `${expected_type}` |
| 3  |  `${actual_type}`  |

## Road Map

- avoiding same value in some field
    - ex case : `old_password` and `new_password`
- get from urls params http
- validation for `base64`
- handle file size on multipart
- extension for generate OpenAPI Spec that support with this package
- image resolution validation
- OR validation (ex : IPv4: true, UUID: true, MultipleCondition: true)
- Custom Validation


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

### Example 2 ( Nested Object Validation )
```go
filterRole := map[string]map_validator.Rules{
    "search":          {Type: reflect.String, Null: true},
    "organization_id": {UUID: true, Null: true},
}
jsonDataRoles := map_validator.NewValidateBuilder().StrictKeys().SetRules(map[string]map_validator.Rules{
    "filter":        {Object: &filterRole, Null: true},
    "rows_per_page": {Type: reflect.Int64, Null: true},
    "page_index":    {Type: reflect.Int64, Null: true},
    "sort": {
        Null:   true,
        IfNull: "FULL_NAME:DESC",
        Type:   reflect.String, Enum: &map_validator.EnumField[any]{
            Items: []string{"FULL_NAME:DESC", "FULL_NAME:ASC", "EMAIL:ASC", "EMAIL:DESC"},
        },
    },
})
```
![image](https://github.com/Rhyanz46/go-map-validator/assets/24217568/9f58dde4-b175-4a4f-9369-fa0974c25942)


### Example 3 ( Echo Framework )
```go
func handleLogin(c echo.Context) error {
    jsonHttp, err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
        "email":    {Email: true, Max: map_validator.SetTotal(100)},
        "password": {Type: reflect.String, Min: map_validator.SetTotal(6), Max: map_validator.SetTotal(30)},
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

### Example 4 ( Bind To Struct )
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


### Example 5 ( Custom message )
```go
payload := map[string]interface{}{"total": 12, "unit": "KG"}
validRole := map[string]map_validator.Rules{
    "total": {
        Type: reflect.Int,
        CustomMsg: map_validator.CustomMsg{
            OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
        },
    },
}
check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
_, err = check.RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
```


### Example 6 ( Regex validator )
```go
payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
validRole := map[string]map_validator.Rules{
    "hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
    "email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
}
check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
_, err = check.RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
```
