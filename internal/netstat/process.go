package netstat

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	TcpType  = "tcp"
	Tcp6Type = "tcp6"
	UdpType  = "udp"
	Udp6Type = "udp6"

	EstablishedState = "ESTABLISHED"
	SynSentState     = "SYN_SENT"
	SynRecvState     = "SYN_RECV"
	FinWait1State    = "FIN_WAIT1"
	FinWait2State    = "FIN_WAIT2"
	TimeWaitState    = "TIME_WAIT"
	CloseState       = "CLOSE"
	CloseWaitState   = "CLOSE_WAIT"
	LastAckState     = "LAST_ACK"
	ListenState      = "LISTEN"
	ClosingState     = "CLOSING"

	FileDescriptors = "/proc/[0-9]*/fd/[0-9]*"

	UnknownUser = "Unknown"
)

var State = map[string]string{
	"01": EstablishedState,
	"02": SynSentState,
	"03": SynRecvState,
	"04": FinWait1State,
	"05": FinWait2State,
	"06": TimeWaitState,
	"07": CloseState,
	"08": CloseWaitState,
	"09": LastAckState,
	"0A": ListenState,
	"0B": ClosingState,
}

var AllowedTypes = []string{TcpType, Tcp6Type, UdpType, Udp6Type}

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
	return netstat(TcpType)
}

func Udp() (Processes, error) {
	return netstat(UdpType)
}

func Tcp6() (Processes, error) {
	return netstat(Tcp6Type)
}

func Udp6() (Processes, error) {
	return netstat(Udp6Type)
}

// Require root access to get information about some processes.
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

	if !contains(AllowedTypes, t) {
		return nil, fmt.Errorf("type %s not allowed", t)
	}

	filename := fmt.Sprintf("/proc/net/%s", t)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	// Return lines without Header line and blank line on the end
	return lines[1 : len(lines)-1], nil

}

func processNetstatLine(line string, fileDescriptors *[]iNode, output chan<- Process) {

	lineArray := removeEmpty(strings.Split(strings.TrimSpace(line), " "))

	ipPortData := strings.Split(lineArray[1], ":")
	ip := convertIp(ipPortData[0])
	port := hexToDec(ipPortData[1])

	foreignData := strings.Split(lineArray[2], ":")
	foreignIP := convertIp(foreignData[0])
	foreignPort := hexToDec(foreignData[1])

	state := State[lineArray[3]]
	uid := getUser(lineArray[7])
	pid := findPid(lineArray[9], fileDescriptors)
	exe := getProcessExe(pid)
	name := getProcessName(exe)

	output <- Process{uid, name, pid, exe, state, ip, int(port), foreignIP, int(foreignPort)}
}

func hexToDec(h string) int64 {
	d, _ := strconv.ParseInt(h, 16, 32)
	return d
}

// Converts the IPv4 to decimal. Have to rearrange the IP because the
// default value is in little Endian order.
func convertIp(ip string) string {

	// Check IP size if greater than 8 is an IPv6 type
	if len(ip) > 8 {
		i := []string{
			ip[30:32],
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
			ip[0:2],
		}
		return fmt.Sprintf(
			"%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v",
			i[14], i[15],
			i[12], i[13],
			i[10], i[11],
			i[8], i[9],
			i[6], i[7],
			i[4], i[5],
			i[2], i[3],
			i[0], i[1],
		)
	}

	i := []int64{hexToDec(ip[6:8]), hexToDec(ip[4:6]), hexToDec(ip[2:4]), hexToDec(ip[0:2])}
	return fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])

}

func findPid(inode string, inodes *[]iNode) string {
	// Loop through all fd dirs of process on /proc to compare the inode and
	// get the pid.

	pid := ""
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
	return cases.Title(language.English).String(name)
}

func getUser(uid string) string {
	u, err := user.LookupId(uid)
	if err != nil {
		return UnknownUser
	}
	return u.Username
}

func removeEmpty(array []string) []string {
	// remove empty data from line
	var newArray []string
	for _, i := range array {
		if i != "" {
			newArray = append(newArray, i)
		}
	}
	return newArray
}

func getInodes() []iNode {
	fileDescriptors, err := getFileDescriptors(FileDescriptors)
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
