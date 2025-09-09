package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func main() {
	// Print header
	fmt.Printf("%s=== hoi-ola System Monitor ===%s\n\n", Cyan, Reset)

	// Get current time
	currentTime := time.Now().Format("15:04:05 MST")

	// Get RAM usage
	ramUsage := getRAMUsage()

	// Get CPU temperature
	cpuTemp := getCPUTemperature()

	// Get GPU temperature
	gpuTemp := getGPUTemperature()

	// Get network speed
	networkSpeed := getNetworkSpeed()

	// Display information with colors
	fmt.Printf("%s[%s]%s\n", Purple, currentTime, Reset)
	fmt.Printf("%sRAM Usage:%s     %.2f%%\n", Green, Reset, ramUsage)
	fmt.Printf("%sCPU Temp:%s      %.2f°C\n", Yellow, Reset, cpuTemp)
	
	// Handle GPU temperature display
	if gpuTemp >= 0 {
		fmt.Printf("%sGPU Temp:%s      %.2f°C\n", Blue, Reset, gpuTemp)
	} else {
		fmt.Printf("%sGPU Temp:%s      Not available\n", Blue, Reset)
	}
	
	fmt.Printf("%sNetwork:%s       RX: %s TX: %s\n", Red, Reset, networkSpeed.RX, networkSpeed.TX)
	fmt.Println()
}

// getRAMUsage returns the RAM usage percentage
func getRAMUsage() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0.0
	}
	return vmStat.UsedPercent
}

// getCPUTemperature returns the CPU temperature in Celsius
func getCPUTemperature() float64 {
	// Try to read from coretemp hwmon (Intel CPUs)
	tempFiles := []string{
		"/sys/class/hwmon/hwmon5/temp1_input", // This is typically the package temperature
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
		"/sys/class/thermal/thermal_zone2/temp",
		"/sys/class/thermal/thermal_zone3/temp",
		"/sys/class/thermal/thermal_zone4/temp",
		"/sys/class/thermal/thermal_zone5/temp",
		"/sys/class/thermal/thermal_zone6/temp",
		"/sys/class/thermal/thermal_zone7/temp",
		"/sys/class/thermal/thermal_zone8/temp",
		"/sys/class/thermal/thermal_zone9/temp",
	}

	for _, file := range tempFiles {
		if temp := readTemperatureFromFile(file); temp > 0 {
			return temp
		}
	}

	// If we can't read from sysfs, try using the 'sensors' command if available
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Look for CPU temperature patterns
		if strings.Contains(line, "CPU") || strings.Contains(line, "Package") || strings.Contains(line, "Tdie") || strings.Contains(line, "Core") {
			// Extract temperature value
			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.Contains(field, "°C") {
					// Remove the °C and + signs
					tempStr := strings.ReplaceAll(field, "°C", "")
					tempStr = strings.ReplaceAll(tempStr, "+", "")
					if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
						return temp
					}
				}
			}
		}
	}

	return 0.0
}

// readTemperatureFromFile reads temperature from a file and converts it to Celsius
func readTemperatureFromFile(filename string) float64 {
	file, err := os.Open(filename)
	if err != nil {
		return 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		tempStr := strings.TrimSpace(scanner.Text())
		temp, err := strconv.Atoi(tempStr)
		if err != nil {
			return 0.0
		}

		// Convert from millidegree Celsius to degree Celsius
		tempFloat := float64(temp) / 1000.0

		// Check if this is a reasonable CPU temperature (between 10°C and 100°C)
		if tempFloat > 10.0 && tempFloat < 100.0 {
			return tempFloat
		}
	}

	return 0.0
}

// getGPUTemperature returns the GPU temperature in Celsius
func getGPUTemperature() float64 {
	// Try nvidia-smi for NVIDIA GPUs
	cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			tempStr := strings.TrimSpace(lines[0])
			if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
				return temp
			}
		}
	}

	// Try reading from AMD GPU sensors (if available)
	amdTempFiles := []string{
		"/sys/class/drm/card0/device/hwmon/hwmon2/temp1_input",
		"/sys/class/drm/card1/device/hwmon/hwmon3/temp1_input",
	}

	for _, file := range amdTempFiles {
		if temp := readTemperatureFromFile(file); temp > 0 {
			return temp
		}
	}

	// If no GPU temperature is available, return -1 to indicate it's not available
	return -1.0
}

// NetworkSpeed represents the network speed
type NetworkSpeed struct {
	RX string
	TX string
}

// getNetworkSpeed returns the current network speed
func getNetworkSpeed() NetworkSpeed {
	// Get initial counters
	initialNetIO, err := net.IOCounters(true)
	if err != nil {
		return NetworkSpeed{RX: "N/A", TX: "N/A"}
	}

	// Wait a short time to get a more accurate measurement
	time.Sleep(500 * time.Millisecond)

	// Get final counters
	finalNetIO, err := net.IOCounters(true)
	if err != nil {
		return NetworkSpeed{RX: "N/A", TX: "N/A"}
	}

	initialTimestamp := time.Now()
	time.Sleep(500 * time.Millisecond)
	finalTimestamp := time.Now()
	
	timeDiff := finalTimestamp.Sub(initialTimestamp).Seconds()

	if timeDiff <= 0 {
		return NetworkSpeed{RX: "0 KB/s", TX: "0 KB/s"}
	}

	var totalRX, totalTX float64

	// Create a map of initial counters for easy lookup
	initialMap := make(map[string]net.IOCountersStat)
	for _, counter := range initialNetIO {
		initialMap[counter.Name] = counter
	}

	// Calculate the difference for each interface
	for _, finalCounter := range finalNetIO {
		// Skip loopback interface
		if finalCounter.Name == "lo" {
			continue
		}

		// Skip interfaces with no initial data
		initialCounter, exists := initialMap[finalCounter.Name]
		if !exists {
			continue
		}

		// Calculate bytes transferred during this interval
		rxBytes := float64(finalCounter.BytesRecv - initialCounter.BytesRecv)
		txBytes := float64(finalCounter.BytesSent - initialCounter.BytesSent)

		totalRX += rxBytes
		totalTX += txBytes
	}

	// Convert to KB/s
	rxSpeed := totalRX / timeDiff / 1024.0
	txSpeed := totalTX / timeDiff / 1024.0

	// Format the speed values
	rxStr := formatSpeed(rxSpeed)
	txStr := formatSpeed(txSpeed)

	return NetworkSpeed{
		RX: rxStr,
		TX: txStr,
	}
}

// formatSpeed formats the speed value with appropriate units
func formatSpeed(speed float64) string {
	if speed < 1024 {
		return fmt.Sprintf("%.1f KB/s", speed)
	}
	return fmt.Sprintf("%.1f MB/s", speed/1024.0)
}