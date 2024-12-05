package main

import (
	"time"

	"pulsar_msg_key_stats/src/stats"

	"github.com/sirupsen/logrus"

	"github.com/apache/pulsar-client-go/pulsar"
	pulsar_log "github.com/apache/pulsar-client-go/pulsar/log"
)

func logging(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		//FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalln("Failed parse log level. Reason: ", err)
	}
	logrus.SetLevel(l)
}

func new_client(url string, trust_cert_file string, cert_file string, key_file string, allow_insecure_connection bool) pulsar.Client {
	var client pulsar.Client
	var err error
	var auth pulsar.Authentication

	if len(cert_file) > 0 || len(key_file) > 0 {
		auth = pulsar.NewAuthenticationTLS(cert_file, key_file)
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:                        url,
		TLSAllowInsecureConnection: allow_insecure_connection,
		Authentication:             auth,
		TLSTrustCertsFilePath:      trust_cert_file,
		Logger:                     pulsar_log.NewLoggerWithLogrus(log),
	})

	if err != nil {
		logrus.Errorf("Failed connect to pulsar. Reason: %+v", err)
	}
	return client
}

func main() {
	opt := from_args()
	logging(opt.loglevel)
	logrus.Infof("%+v", opt)

	source_client := new_client(opt.sourcepulsar, opt.sourcetrustcerts, opt.sourcecertfile, opt.sourcekeyfile, opt.sourceallowinsecureconnection)

	defer source_client.Close()

	if opt.pprofon {
		go activate_profiling(opt.pprofdir, time.Duration(opt.pprofduration)*time.Second)
	}

	stats.Monitor_topic(source_client, opt.sourcetopic, opt.sourcetopicnumberpartitions, opt.sourcesubscription, opt.sourcename)
}
