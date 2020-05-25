module fabric-tools

go 1.13

require (
	github.com/fsouza/go-dockerclient v1.6.5 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/hyperledger/fabric v2.1.0+incompatible
	// github.com/hyperledger/fabric v1.4.3
	// github.com/hyperledger/fabric v2.1.0+incompatible
	github.com/hyperledger/fabric-protos-go v0.0.0-20200506201313-25f6564b9ac4 // indirect
	// github.com/hyperledger/fabric v2.1.0+incompatible
	// github.com/hyperledger/fabric-sdk-go v1.0.0-beta1
	github.com/hyperledger/fabric-sdk-go v1.0.0-alpha5.0.20190411180201-5a9a0e749e4f
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190411180201-5a9a0e749e4f //v0.0.0 17677af803a1107753a459f54c226baf6fb3220e
	// github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e // indirect v1.4.3 //v0.0.0-20190411180201-5a9a0e749e4f
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/sykesm/zap-logfmt v0.0.3 // indirect
	github.com/tedsuo/ifrit v0.0.0-20191009134036-9a97d0632f00 // indirect
)

// replace github.com/hyperledger/fabric => github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v1.4.3
