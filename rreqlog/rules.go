package rreqlog

import "github.com/gin-gonic/gin"

type CustomField func(*gin.Context, ...string) string
