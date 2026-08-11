package resp

import (
        "github.com/gin-gonic/gin"
        "net/http"
)

func OK(ctx *gin.Context, msg string, data any, title ...string) {
        resp := gin.H{
                "message": msg,
        }

        if len(title) > 0 {
                resp[title[0]] = msg
                delete(resp, "message")
        }

        if data != nil {
                resp["data"] = data
        }

        ctx.JSON(http.StatusOK, resp)
}

func Error(ctx *gin.Context, status int, err any) {
        ctx.JSON(status, gin.H{
                "error": err,
        })
}
