package main

import (
	"github.com/jnovack/flag"
)

type opt struct {
	sourcepulsar                  string
	sourcetopic                   string
	sourcetopicnumberpartitions   uint
	sourcesubscription            string
	sourcename                    string
	sourcetrustcerts              string
	sourcecertfile                string
	sourcekeyfile                 string
	sourceallowinsecureconnection bool

	pprofon       bool
	pprofdir      string
	pprofduration uint

	prometheusport uint
	loglevel       string
}

func from_args() opt {

	var opt opt

	flag.StringVar(&opt.sourcepulsar, "source_pulsar", "pulsar://localhost:6650", "Source pulsar address")
	flag.StringVar(&opt.sourcetopic, "source_topic", "persistent://public/default/in", "Source topic name")
	flag.UintVar(&opt.sourcetopicnumberpartitions, "source_topic_number_partitions", 1, "Source topic number of partitions")
	flag.StringVar(&opt.sourcesubscription, "source_subscription", "pulsar_msg_key_stats", "Source subscription name")
	flag.StringVar(&opt.sourcename, "source_name", "aggregator_consumer", "Source consumer name")
	flag.StringVar(&opt.sourcetrustcerts, "source_trust_certs", "", "Path for source pem file, for ca.cert")
	flag.StringVar(&opt.sourcecertfile, "source_cert_file", "", "Path for source cert.pem file")
	flag.StringVar(&opt.sourcekeyfile, "source_key_file", "", "Path for source key-pk8.pem file")
	flag.BoolVar(&opt.sourceallowinsecureconnection, "source_allow_insecure_connection", false, "Source allow insecure connection")

	flag.BoolVar(&opt.pprofon, "pprof_on", false, "Profoling on?")
	flag.StringVar(&opt.pprofdir, "pprof_dir", "./pprof", "Directory for pprof file")
	flag.UintVar(&opt.pprofduration, "pprof_duration", 60*4, "Number of seconds to run pprof")

	flag.UintVar(&opt.prometheusport, "prometheus_port", 7700, "Prometheous port")
	flag.StringVar(&opt.loglevel, "log_level", "info", "Logging level: panic - fatal - error - warn - info - debug - trace")

	flag.Parse()

	return opt

}
