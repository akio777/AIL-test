package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var requestLogger = logrus.New()

func Init() *logrus.Logger {
	requestLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	requestLogger.SetLevel(logrus.InfoLevel)
	requestLogger.SetReportCaller(true)
	return requestLogger
}

func RequestLogger(c *fiber.Ctx) error {
	entry := requestLogger.WithFields(logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
		"ip":     c.IP(),
	})
	entry.Info("")

	return c.Next()

}
