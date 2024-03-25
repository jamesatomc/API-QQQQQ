package middleware

import (
	"log"
	"time"


	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
      // ดึง token จาก header ของ request
      token := c.Request.Header.Get("Authorization")
  
      // ตรวจสอบ token ว่าถูกต้องหรือไม่
      // ...
  
      // กำหนดค่า token ให้กับ context
      c.Set("token", token)
  
      // บันทึกเวลาเริ่มต้น
      t := time.Now()
  
      // เรียกใช้ middleware ถัดไป
      c.Next()
  
      // บันทึกเวลาแฝง
      latency := time.Since(t)
  
      // บันทึก log ข้อมูลต่างๆ
      log.Print(latency)
      log.Println(c.Writer.Status())
    }
  }
