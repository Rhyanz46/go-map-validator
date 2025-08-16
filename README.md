# Simple Validator using struct properties as rules

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
- Unique Value
  - ex case : `old_password` and `new_password` cant be using same value
- RequiredWithout check
  - ex case : the field `flavor` is required if `custom_flavor` is null
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
  - on invalid regex message : ✅ ready
  - on type not match message : ✅ ready
  - on min/max data message : ✅ ready
  - on unique values error : ✅ ready
  - on null data message : ⌛ not ready
  - on enum value not match : ⌛ not ready
  - on `RequiredWithout` error : ⌛ not ready

## On Progress

- validation for one data value only

## Custom Message Variables

| No |      Variable Name       |
|:--:|:------------------------:| 
| 1  |        `${field}`        |
| 2  |    `${expected_type}`    |
| 3  |     `${actual_type}`     |
| 4  |    `${actual_length}`    |
| 5  | `${expected_min_length}` |
| 6  | `${expected_max_length}` |
| 7  |    `${unique_origin}`    |
| 8  |    `${unique_target}`    |


## Road Map
- errors detail mode
- get from urls params http
- validation for `base64`
- handle file size on multipart
- extension for generate OpenAPI Spec that support with this package
- image resolution validation
- multi validation on one field (ex : IPv4: true, UUID: true)


## example :

### Example 1
```go
payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}

op, err := map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
    Rules: map[string]map_validator.Rules{
        "jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
        "hoby":          {Type: reflect.String, Null: false},
        "menikah":       {Type: reflect.Bool, Null: false},
    },
}).Load(payload)
if err != nil {
    t.Fatalf("load error: %s", err)
}
if _, err = op.RunValidate(); err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}

op, err = map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
    Rules: map[string]map_validator.Rules{
        "jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
        "hoby":          {Type: reflect.Int, Null: false},
        "menikah":       {Type: reflect.Bool, Null: false},
    },
}).Load(payload)
if err != nil {
    t.Fatalf("load error: %s", err)
}
if _, err = op.RunValidate(); err == nil {
    t.Error("Expected have an error, but you got no error")
}
```

### Example 2 ( Nested Object Validation )
```go
filterRules := map_validator.BuildRoles().
  SetRule("search", map_validator.Rules{Type: reflect.String, Null: true}).
  SetRule("organization_id", map_validator.Rules{UUID: true, Null: true}).
  Done()

parent := map_validator.BuildRoles().
  SetRule("filter", map_validator.Rules{Object: &filterRules, Null: true}).
  SetRule("rows_per_page", map_validator.Rules{Type: reflect.Int64, Null: true}).
  SetRule("page_index", map_validator.Rules{Type: reflect.Int64, Null: true}).
  SetRule("sort", map_validator.Rules{
      Null:   true,
      IfNull: "FULL_NAME:DESC",
      Type:   reflect.String,
      Enum:   &map_validator.EnumField[any]{Items: []string{"FULL_NAME:DESC", "FULL_NAME:ASC", "EMAIL:ASC", "EMAIL:DESC"}},
  }).
  SetSetting(map_validator.Setting{Strict: true}).
  Done()

op, err := map_validator.NewValidateBuilder().SetRules(parent).Load(map[string]interface{}{})
if err != nil { t.Fatal(err) }
_, _ = op.RunValidate()
```
![image](https://github.com/Rhyanz46/go-map-validator/assets/24217568/9f58dde4-b175-4a4f-9369-fa0974c25942)


### Example 3 ( Echo Framework )
```go
func handleLogin(c echo.Context) error {
    jsonHttp, err := map_validator.NewValidateBuilder().SetRules(
        map_validator.BuildRoles().
            SetRule("email", map_validator.Rules{Email: true, Max: map_validator.SetTotal(100)}).
            SetRule("password", map_validator.Rules{Type: reflect.String, Min: map_validator.SetTotal(6), Max: map_validator.SetTotal(30)}).
            Done(),
    ).LoadJsonHttp(c.Request())
    if err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    if _, err := op.RunValidate(); err != nil {
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
    JK      string `json:"jenis_kelamin"`
    Hoby    string `json:"hoby"`
    Menikah bool   `json:"menikah"`
}

payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
err := map_validator.NewValidateBuilder().SetRules(
    map_validator.BuildRoles().
        SetRule("jenis_kelamin", map_validator.Rules{Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}}).
        SetRule("hoby", map_validator.Rules{Type: reflect.String, Null: false}).
        SetRule("menikah", map_validator.Rules{Type: reflect.Bool, Null: false}).
        Done(),
).Load(payload).RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}

testBind := &Data{}
if testBind.JK != "" {
    t.Errorf("Expected : '' But you got : %s", testBind.JK)
}
if err := extraCheck.Bind(testBind); err != nil { t.Fatal(err) }

if testBind.JK != payload["jenis_kelamin"] {
    t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
}

```


### Example 5 ( Custom message )
```go
payload := map[string]interface{}{"total": 12, "unit": "KG"}
validRole := map_validator.BuildRoles().
    SetRule("total", map_validator.Rules{
        Type: reflect.Int,
        CustomMsg: map_validator.CustomMsg{
            OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
        },
    }).
    Done()

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
payload := map[string]interface{}{ "hp": "+62567888", "email": "dev@ariansaputra.com" }
validRole := map_validator.BuildRoles().
    SetRule("hp", map_validator.Rules{ RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$` }).
    SetRule("email", map_validator.Rules{ RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` }).
    Done()
check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
_, err = check.RunValidate()
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
}
```

### Example 7 ( Unique value )
```go
role := map_validator.BuildRoles().
    SetRule("password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true}).
    SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true}).
    Done()
payload := map[string]interface{}{
    "password":     "sabalong",
    "new_password": "sabalong",
}
check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
if err != nil {
    t.Errorf("Expected not have error, but got error : %s", err)
    return
}
expected := "value of 'password' and 'new_password' fields must be different"
_, err = check.RunValidate()
if err.Error() != expected {
    t.Errorf("Expected :%s. But you got : %s", expected, err)
}
```

### Example 8 ( HTTP JSON with ListObject + Bind )
```go
// Define your request struct with a slice field
type GoodsRequest struct {
    Name        string  `json:"name"`
    Weight      float64 `json:"weight"`
    Quantity    int     `json:"quantity"`
    Description string  `json:"description"`
}
type CreateOrderRequest struct {
    SenderID                string         `json:"sender_id"`
    SenderAddress           string         `json:"sender_address"`
    SenderAddressCity       string         `json:"sender_address_city"`
    SenderAddressProvince   string         `json:"sender_address_province"`
    SenderLatitude          float64        `json:"sender_latitude"`
    SenderLongitude         float64        `json:"sender_longitude"`
    ReceiverName            string         `json:"receiver_name"`
    ReceiverPhone           string         `json:"receiver_phone"`
    ReceiverAddress         string         `json:"receiver_address"`
    ReceiverAddressCity     string         `json:"receiver_address_city"`
    ReceiverAddressProvince string         `json:"receiver_address_province"`
    ReceiverLatitude        float64        `json:"receiver_latitude"`
    ReceiverLongitude       float64        `json:"receiver_longitude"`
    Note                    string         `json:"note"`
    Goods                   []GoodsRequest `json:"goods"`
}

// Sample JSON (goods as an array)
jsonStr := `{
  "goods": [
    {"description":"string","name":"string","quantity":1,"weight":0}
  ],
  "note":"string",
  "receiver_address":"string",
  "receiver_address_city":"string",
  "receiver_address_province":"string",
  "receiver_latitude":0,
  "receiver_longitude":0,
  "receiver_name":"string",
  "receiver_phone":"string",
  "sender_address":"string",
  "sender_address_city":"string",
  "sender_address_province":"string",
  "sender_id":"8aa3e797-2453-442f-b1d0-50f7d815bcaf",
  "sender_latitude":0,
  "sender_longitude":0
}`

req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(jsonStr))
req.Header.Set("Content-Type", "application/json")

rules := map_validator.BuildRoles().
  SetRule("sender_id", map_validator.Rules{Type: reflect.String, UUID: true}).
  SetRule("sender_address", map_validator.Rules{Type: reflect.String}).
  SetRule("sender_address_city", map_validator.Rules{Type: reflect.String}).
  SetRule("sender_address_province", map_validator.Rules{Type: reflect.String}).
  SetRule("sender_latitude", map_validator.Rules{Type: reflect.Float64}).
  SetRule("sender_longitude", map_validator.Rules{Type: reflect.Float64}).
  SetRule("receiver_name", map_validator.Rules{Type: reflect.String}).
  SetRule("receiver_phone", map_validator.Rules{Type: reflect.String}).
  SetRule("receiver_address", map_validator.Rules{Type: reflect.String}).
  SetRule("receiver_address_city", map_validator.Rules{Type: reflect.String}).
  SetRule("receiver_address_province", map_validator.Rules{Type: reflect.String}).
  SetRule("receiver_latitude", map_validator.Rules{Type: reflect.Float64}).
  SetRule("receiver_longitude", map_validator.Rules{Type: reflect.Float64}).
  SetRule("note", map_validator.Rules{Type: reflect.String}).
  SetRule("goods", map_validator.Rules{ListObject: map_validator.BuildRoles().
    SetRule("name", map_validator.Rules{Type: reflect.String}).
    SetRule("weight", map_validator.Rules{Type: reflect.Float64}).
    SetRule("quantity", map_validator.Rules{Type: reflect.Int, Min: map_validator.SetTotal(1)}).
    SetRule("description", map_validator.Rules{Type: reflect.String}),
  }).
  Done()

jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
if err != nil { panic(err) }

extra, err := jsonHttp.RunValidate()
if err != nil { panic(err) }

var reqDTO CreateOrderRequest
if err := extra.Bind(&reqDTO); err != nil { panic(err) }

// Note: "goods" MUST be an array when using ListObject.
// If you send an object instead, the validator returns:
//   "field 'goods' is not valid list object"
```
