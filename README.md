# hoi-ola

A simple CLI tool for monitoring system resources with colorful output.

## Features

- Shows RAM usage percentage
- Displays CPU temperature
- Displays GPU temperature (NVIDIA/AMD)
- Shows network speed (RX/TX)
- Colorful output with timestamp

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/vedantkalkar19/Hoi-Ola.git
   cd hoi-ola
   ```

2. Build the application:
   ```
   go build -o hoi-ola
   ```

3. Run the application:
   ```
   ./hoi-ola
   ```

## Requirements

- Go 1.16 or higher
- For GPU temperature monitoring:
  - NVIDIA GPUs: `nvidia-smi` must be installed
  - AMD GPUs: Appropriate drivers and sysfs support

## Usage

Simply run the binary to get a snapshot of your system's current status:

```
./hoi-ola
```

The application will display the current system time, RAM usage, CPU temperature, GPU temperature, and network speeds, then exit.

## Color Codes

- **System Time**: Purple
- **RAM Usage**: Green
- **CPU Temperature**: Yellow
- **GPU Temperature**: Blue
- **Network Speed**: Red

## Example Output

```
=== hoi-ola System Monitor ===

[16:09:43 IST]
RAM Usage:     51.01%
CPU Temp:      72.00°C
GPU Temp:      62.00°C
Network:       RX: 1.9 KB/s TX: 0.7 KB/s
```

## License

MIT
