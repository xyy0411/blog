package flag

import "flag"

type Options struct {
	File    string
	DB      bool
	Version bool
}

var flagOptions = new(Options)

func Parse() {
	flag.StringVar(&flagOptions.File, "f", "/config.config.yaml", "配置文件")
	flag.BoolVar(&flagOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&flagOptions.Version, "v", false, "版本")
	flag.Parse()
}
