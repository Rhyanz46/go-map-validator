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
