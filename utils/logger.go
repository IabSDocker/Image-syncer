package utils

import (
    "github.com/sirupsen/logrus"
    "time"
)

type CSTFormatter struct {
    logrus.TextFormatter
}

func (f *CSTFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    // Convert time to East 8 timezone (UTC+8)
    localTime := entry.Time.In(time.FixedZone("CST", 8*3600))
    entry.Time = localTime

    return f.TextFormatter.Format(entry)
}

func CreateLogger() *logrus.Logger {
    logger := logrus.New()

    // Set custom formatter with CST (UTC+8) timezone
    logger.SetFormatter(&CSTFormatter{
        TextFormatter: logrus.TextFormatter{
            TimestampFormat: time.RFC3339,
            FullTimestamp:   true,
        },
    })

    return logger
}