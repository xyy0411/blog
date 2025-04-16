package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xyy0411/blog/models"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FormatMatchingInfo 将 models.Matching 结构体中的信息格式化为字符串
func FormatMatchingInfo(value1 models.Matching) string {

	var result strings.Builder
	result.WriteString("匹配成功啦!对手信息如下:\n名称:")
	result.WriteString(value1.UserName)
	result.WriteString("\n群号:")
	result.WriteString(strconv.FormatInt(value1.GroupID, 10))
	result.WriteString("\nQQ号:")
	result.WriteString(strconv.FormatInt(value1.UserID, 10))
	result.WriteString("\n支持的软件:\n")
	for _, s := range value1.OnlineSoftware {
		var t string
		switch s.Type {
		case 0:
			t = "主副皆可"
		case 1:
			t = "主机"
		case 2:
			t = "副机"
		}
		result.WriteString("软件:")
		result.WriteString(s.Name)
		result.WriteString("   主副机:")
		result.WriteString(t)
		result.WriteString("\n")
	}
	return result.String()
}

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
