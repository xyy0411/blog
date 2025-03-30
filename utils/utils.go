package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"time"
)

// StructToMap 将结构体转map 忽视字段请使用tag structToMap:"ignore"
func StructToMap(obj any) map[string]any {
	result := make(map[string]any)
	if obj == nil {
		return result
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return result
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		if tag := fieldType.Tag.Get("structToMap"); tag == "ignore" {
			continue
		}
		fieldName := fieldType.Name
		result[fieldName] = field.Interface()
	}
	return result
}

func GenerateJWT(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("secret"))
}

func CheckPassword(pwd string, hash string) bool {
	return pwd == hash
}

func ParseJWT(tokenstring string) (string, error) {
	if len(tokenstring) > 7 && tokenstring[:7] == "Bearer " {
		tokenstring = tokenstring[7:]
	}

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", errors.New("解析token错误")
	}

	return claims["name"].(string), nil
}
