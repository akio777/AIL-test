package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestLogger(c *fiber.Ctx) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	entry := logrus.WithFields(logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
		"ip":     c.IP(),
	})
	entry.Info("")

	return c.Next()

}
