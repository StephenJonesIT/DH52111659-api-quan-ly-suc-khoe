package utils

func CalculateBMI(weight, height int) float64 {
	if height <= 0 {
		return 0.0 // Avoid division by zero
	}
	heightInMeters := float64(height) / 100.0 // Convert height from cm to meters
	bmi := float64(weight) / (heightInMeters * heightInMeters)
	return bmi
}