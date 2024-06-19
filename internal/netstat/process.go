package netstat

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	TCP  = "tcp"
	TCP6 = "tcp6"
	UDP  = "udp"
	UDP6 = "udp6"

	PROC_TCP  = "/proc/net/tcp"
	PROC_UDP  = "/proc/net/udp"
	PROC_TCP6 = "/proc/net/tcp6"
	PROC_UDP6 = "/proc/net/udp6"

	ESTABLISHED_STATE = "ESTABLISHED"
	SYN_SENT_STATE    = "SYN_SENT"
	SYN_RECV_STATE    = "SYN_RECV"
	FIN_WAIT1_STATE   = "FIN_WAIT1"
	FIN_WAIT2_STATE   = "FIN_WAIT2"
	TIME_WAIT_STATE   = "TIME_WAIT"
	CLOSE_STATE       = "CLOSE"
	CLOSE_WAIT_STATE  = "CLOSE_WAIT"
	LAST_ACK_STATE    = "LAST_ACK"
	LISTEN_STATE      = "LISTEN"
	CLOSING_STATE     = "CLOSING"
)

var STATE = map[string]string{
	"01": ESTABLISHED_STATE,
	"02": SYN_SENT_STATE,
	"03": SYN_RECV_STATE,
	"04": FIN_WAIT1_STATE,
	"05": FIN_WAIT2_STATE,
	"06": TIME_WAIT_STATE,
	"07": CLOSE_STATE,
	"08": CLOSE_WAIT_STATE,
	"09": LAST_ACK_STATE,
	"0A": LISTEN_STATE,
	"0B": CLOSING_STATE,
}

type Process struct {
	User        string
	Name        string
	Pid         string
	Exe         string
	State       string
	Ip          string
	Port        int
	ForeignIp   string
	ForeignPort int
}

type Processes []Process

type iNode struct {
	path string
	link string
}

func Tcp() (Processes, error) {
	return netstat(TCP)
}

func Udp() (Processes, error) {
	return netstat(PROC_UDP6)
}

func Tcp6() (Processes, error) {
	return netstat("tcp6")
}

func Udp6() (Processes, error) {
	return netstat("udp6")
}

// Require root acess to get information about some processes.
func netstat(t string) (Processes, error) {

	data, err := getData(t)
	if err != nil {
		return Processes{}, nil
	}

	processes := make([]Process, len(data))
	res := make(chan Process, len(data))

	inodes := getInodes()
	for _, line := range data {
		go processNetstatLine(line, &inodes, res)
	}

	for i := range data {
		p := <-res
		processes[i] = p
	}

	return processes, nil
}

func getData(t string) ([]string, error) {

	var proc_t string

	switch t {
	case TCP:
		proc_t = PROC_TCP
	case UDP:
		proc_t = PROC_UDP
	case TCP6:
		proc_t = PROC_TCP6
	case UDP6:
		proc_t = PROC_UDP6
	default:
		return nil, fmt.Errorf("invalid type %s", t)
	}

	data, err := os.ReadFile(proc_t)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	// Return lines without Header line and blank line on the end
	return lines[1 : len(lines)-1], nil

}

func hexToDec(h string) int64 {
	d, _ := strconv.ParseInt(h, 16, 32)
	return d
}

// Converts the ipv4 to decimal. Have to rearrange the ip because the
// default value is in little Endian order.
func convertIp(ip string) string {

	// Check ip size if greater than 8 is a ipv6 type
	if len(ip) > 8 {
		i := []string{ip[30:32],
			ip[28:30],
			ip[26:28],
			ip[24:26],
			ip[22:24],
			ip[20:22],
			ip[18:20],
			ip[16:18],
			ip[14:16],
			ip[12:14],
			ip[10:12],
			ip[8:10],
			ip[6:8],
			ip[4:6],
			ip[2:4],
			ip[0:2]}
		return fmt.Sprintf(
			"%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v",
			i[14], i[15], i[13], i[12],
			i[10], i[11], i[8], i[9],
			i[6], i[7], i[4], i[5],
			i[2], i[3], i[0], i[1],
		)
	}

	i := []int64{hexToDec(ip[6:8]), hexToDec(ip[4:6]), hexToDec(ip[2:4]), hexToDec(ip[0:2])}
	return fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])

}

func findPid(inode string, inodes *[]iNode) string {
	// Loop through all fd dirs of process on /proc to compare the inode and
	// get the pid.

	pid := "-"

	re := regexp.MustCompile(inode)
	for _, item := range *inodes {
		out := re.FindString(item.link)
		if len(out) != 0 {
			pid = strings.Split(item.path, "/")[2]
		}
	}
	return pid
}

func getProcessExe(pid string) string {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

func getProcessName(exe string) string {
	n := strings.Split(exe, "/")
	name := n[len(n)-1]
	return strings.Title(name)
}

func getUser(uid string) string {
	u, err := user.LookupId(uid)
	if err != nil {
		return "Unknown"
	}
	return u.Username
}

func removeEmpty(array []string) []string {
	// remove empty data from line
	var new_array []string
	for _, i := range array {
		if i != "" {
			new_array = append(new_array, i)
		}
	}
	return new_array
}

func processNetstatLine(line string, fileDescriptors *[]iNode, output chan<- Process) {
	line_array := removeEmpty(strings.Split(strings.TrimSpace(line), " "))
	ip_port := strings.Split(line_array[1], ":")
	ip := convertIp(ip_port[0])
	port := hexToDec(ip_port[1])

	// foreign ip and port
	fip_port := strings.Split(line_array[2], ":")
	fip := convertIp(fip_port[0])
	fport := hexToDec(fip_port[1])

	state := STATE[line_array[3]]
	uid := getUser(line_array[7])
	pid := findPid(line_array[9], fileDescriptors)
	exe := getProcessExe(pid)
	name := getProcessName(exe)
	output <- Process{uid, name, pid, exe, state, ip, int(port), fip, int(fport)}
}

func getInodes() []iNode {
	fileDescriptors, err := getFileDescriptors("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		return nil
	}

	inodes := make([]iNode, len(fileDescriptors))
	res := make(chan iNode, len(fileDescriptors))

	go func(fileDescriptors *[]string, output chan<- iNode) {
		for _, item := range *fileDescriptors {
			link, _ := os.Readlink(item)
			output <- iNode{item, link}
		}
	}(&fileDescriptors, res)

	for range fileDescriptors {
		inode := <-res
		inodes = append(inodes, inode)
	}

	return inodes
}

func getFileDescriptors(pattern string) ([]string, error) {
	d, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return d, nil
}
