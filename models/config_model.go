package models

type Run struct {
	ID   int    `mapstructure:"id"`
	Name string `mapstructure:"name"`
	Ipfs string `mapstructure:"ipfs"`
	Rest string `mapstructure:"rest"`
}

type RunList struct {
	Runs []Run `mapstructure:"runs"`
}
