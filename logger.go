package logger

import (
	"net"
	"os"

	"github.com/michaelvlaar/logrusly"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/shiena/ansicolor"
)

const addressWasNotFound = "suitable address was not found"

func ipAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.Wrap(err, "failed to")
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return addressWasNotFound, nil
}

// New creates new logger instance with hook to send logs to loggly
func New(logglyToken, logglyHost, logglyTag string, logLevel string) (logrus.FieldLogger, error) {
	// Setup logging
	logr := logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse logLevel")
	}
	logr.SetLevel(level)

	// Setting hook to send logs to loggly
	hook := logrusly.NewLogglyHook(logglyToken, logglyHost, level, logglyTag)
	logr.Hooks.Add(hook)

	logr.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	logr.Out = ansicolor.NewAnsiColorWriter(os.Stdout)

	addr, err := ipAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ip address")
	}
	// Adding default field with node ip
	logglyLogger := logr.WithFields(logrus.Fields{
		"node": addr,
	})
	return logglyLogger, nil
}
