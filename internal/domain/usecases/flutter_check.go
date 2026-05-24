package usecases

import (
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/flutter"
)

type FlutterCheckResult struct {
	Available bool
	Version   string
	Doctor    flutter.DoctorResult
}

func CheckFlutter() FlutterCheckResult {
	result := flutter.RunDoctorCheck()
	return FlutterCheckResult{
		Available: result.FlutterInstalled,
		Version:   result.FlutterVersion,
		Doctor:    result,
	}
}
