module github.com/feitianlove/worker

go 1.13

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/feitianlove/golib v0.0.0-20210715152206-b5bed50cfe70
	github.com/feitianlove/web v0.0.0-20210718154758-45ad710aa3b7
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
