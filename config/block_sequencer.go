package config

type BlockSequencerConf struct {
	Active bool
}

var defaultBlockSequencerConf = &BlockSequencerConf{
	Active: true,
}
