package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func UnmarshalFormValue(ctx *gin.Context, key string, v interface{}) error{
	jsonValue := ctx.PostForm(key)
	if jsonValue == ""{
		return nil
	}

	if err := json.Unmarshal([]byte(jsonValue), v); err != nil {
		return err
	}

	return nil
}