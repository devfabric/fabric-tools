package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"code.uni-ledger.com/switch/fabric-tools/types"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Cfg 提供配置
type Cfg struct {
	core.ConfigProvider
	ChainCodes types.ChainCodes
	Opts
}

// NewCfg 返回配置实体
func NewCfg(opts *Opts) (*Cfg, error) {
	// 读取配置文件
	sdkcfg := config.FromFile(opts.SdkCfgFile)
	provider, err := sdkcfg()
	if err != nil {
		fmt.Println("Get fabric configuration item from config file err:", err)
		os.Exit(1)
	}

	// 设置日志
	if opts.LogLevel == "" {
		lvl, _ := provider[0].Lookup("client.logging.level")
		lvlstr, ok := lvl.(string)
		if ok {
			opts.LogLevel = lvlstr
		}
	}
	if opts.Logger == nil {
		opts.Logger = logrus.New()
		lvl, err := logrus.ParseLevel(opts.LogLevel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		opts.Logger.SetLevel(lvl)
		opts.Logger.SetReportCaller(true)
		opts.Logger.SetFormatter(newTextFormatter())
	}

	cfg := Cfg{
		ConfigProvider: sdkcfg,
		Opts:           *opts,
	}

	cfg.UserCfg()

	return &cfg, nil
}

// UserCfg 获取自定义配置
func (C *Cfg) UserCfg() {
	cfgFile := C.Opts.SdkCfgFile
	path, file := filepath.Split(cfgFile)
	file = strings.TrimSuffix(file, ".yaml")

	cfg := viper.New()
	cfg.SetConfigName(file)   // 配置文件名（不带后缀）
	cfg.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	cfg.AddConfigPath(path)   // 可以在工作目录下查找
	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Println("Get user configuration item from config file err:", err)
		os.Exit(1)
	}

	cc := types.ChainCodes{}
	err = cfg.UnmarshalKey("chaincode", &cc)
	if err != nil {
		fmt.Printf("Get chaincode configuration item from %s err:%v\n", C.Opts.SdkCfgFile, err)
		os.Exit(1)
	}
	C.ChainCodes = cc
}

// Provider 返回配置配置文件实体
func (C *Cfg) Provider() core.ConfigProvider {
	return C.ConfigProvider
}

// NewTextFormatter 日志格式
func newTextFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		DisableColors:          false,                     //是否输出颜色
		DisableTimestamp:       false,                     //是否禁用时间戳
		FullTimestamp:          true,                      //true：时间，false:序列号
		TimestampFormat:        "2006-01-02 15:04:05.999", //时间戳格式
		DisableLevelTruncation: false,                     //4字日志等级
		CallerPrettyfier:       callerTextPrettyfierFunc,  // 设置ReportCaller(调用信息)格式
	}
}

//callerTextPrettyfierFunc 设置出格式：输出文件行号，不输出函数
func callerTextPrettyfierFunc(frame *runtime.Frame) (function string, file string) {
	return "", fmt.Sprintf("%s.%d", path.Base(frame.File), frame.Line)
}
