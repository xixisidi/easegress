package option

import (
	"flag"
	"os"
	"strings"
	"time"

	"common"
)

var (
	// cluster stuff
	ClusterGroup          string
	MemberMode            string
	MemberName            string
	Peers                 []string
	OPLogMaxSeqGapToPull  uint16
	OPLogPullMaxCountOnce uint16
	OPLogPullInterval     time.Duration
	OPLogPullTimeout      time.Duration

	Host                           string
	CertFile, KeyFile              string
	Stage                          string
	ConfigHome, LogHome            string
	CpuProfileFile, MemProfileFile string

	ShowVersion bool
)

func init() {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "node0"
	}

	clusterGroup := flag.String("group", "default", "specify cluster group")
	memberMode := flag.String("mode", "read", "specify member mode (read or write)")
	memberName := flag.String("name", hostName, "specify member name")
	peers := flag.String("peers", "", "specify address list of peer members (seprated by comma)")
	opLogMaxSeqGapToPull := new(uint16)
	flag.Var(common.NewUint16Value(5, opLogMaxSeqGapToPull), "oplog_max_seq_gap_to_pull",
		"specify max gap of sequnce of operation logs deciding whether to wait for missing operations or not")
	opLogPullMaxCountOnce := new(uint16)
	flag.Var(common.NewUint16Value(5, opLogPullMaxCountOnce), "oplog_pull_max_count_once",
		"specify max count of pulling operation logs once")
	opLogPullInterval := new(uint16)
	flag.Var(common.NewUint16Value(10, opLogPullInterval), "oplog_pull_interval",
		"specify interval of pulling operation logs in second")
	opLogPullTimeout := new(uint16)
	flag.Var(common.NewUint16Value(30, opLogPullTimeout), "oplog_pull_timeout",
		"specify timeout of pulling operation logs in second")

	host := flag.String("host", "localhost", "specify listen host")
	certFile := flag.String("certfile", "", "specify cert file, "+
		"downgrade HTTPS(10443) to HTTP(10080) if it is set empty or inexistent file")
	keyFile := flag.String("keyfile", "", "specify key file, "+
		"downgrade HTTPS(10443) to HTTP(10080) if it is set empty or inexistent file")
	stage := flag.String("stage", "debug", "sepcify runtime stage (debug, test, prod)")
	configHome := flag.String("config", common.CONFIG_HOME_DIR, "sepcify config home path")
	logHome := flag.String("log", common.LOG_HOME_DIR, "specify log home path")
	cpuProfileFile := flag.String("cpuprofile", "", "specify cpu profile output file, "+
		"cpu profiling will be fully disabled if not provided")
	memProfileFile := flag.String("memprofile", "", "specify heap dump file, "+
		"memory profiling will be fully disabled if not provided")
	showVersion := flag.Bool("version", false, "ouptput version information")

	flag.Parse()

	ClusterGroup = *clusterGroup
	MemberMode = *memberMode
	MemberName = *memberName
	OPLogMaxSeqGapToPull = *opLogMaxSeqGapToPull
	OPLogPullMaxCountOnce = *opLogPullMaxCountOnce
	OPLogPullInterval = time.Duration(*opLogPullInterval) * time.Second
	OPLogPullTimeout = time.Duration(*opLogPullTimeout) * time.Second
	Peers = make([]string, 0)
	for _, peer := range strings.Split(*peers, ",") {
		peer = strings.TrimSpace(peer)
		if len(peer) > 0 {
			Peers = append(Peers, peer)
		}
	}

	Host = *host
	CertFile = *certFile
	KeyFile = *keyFile
	Stage = *stage
	ConfigHome = *configHome
	LogHome = *logHome
	CpuProfileFile = *cpuProfileFile
	MemProfileFile = *memProfileFile
	ShowVersion = *showVersion
}