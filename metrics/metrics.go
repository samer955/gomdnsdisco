package metrics

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"os"
	"strconv"
)
import "github.com/shirou/gopsutil/v3/host"

func New() Metrics {
	return Metrics{}
}

type Metrics struct {
	System System
	Cpu    Cpu
}

type System struct {
	Hostname string `json:"hostname"`
	Os       string `json:"os"`
	Platform string `json:"platform"`
	Version  string `json:"version"`
}

type Cpu struct {
	Model string `json:"model"`
	Usage string `json:"usage"`
}

func (s *System) GetSystemInformation() {

	hostname, _ := os.Hostname()
	s.Hostname = hostname

	hostStat, err := host.Info()

	if err != nil {
		s.Os, s.Version, s.Platform = "", "", ""
		return
	}
	s.Os, s.Platform, s.Version = hostStat.OS, hostStat.Platform, hostStat.PlatformVersion

}

func (c *Cpu) GetUsagePercent() {

	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Unable to get Cpu percent Usage")
		c.Usage = ""
		return
	}

	c.Usage = strconv.Itoa(int(percent[0]))

}

func (c *Cpu) GetModel() {

	cpuStat, err := cpu.Info()

	if err != nil {
		c.Model = ""
		return
	}
	c.Model = cpuStat[0].ModelName

}
