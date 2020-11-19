module godbg

go 1.15

//replace github.com/stromland/cobra-prompt => ../../../../Github/cobra-prompt
replace github.com/stromland/cobra-prompt => github.com/hitzhangjie/cobra-prompt v0.0.0-20201118115017-e51f84e8374c

require (
	github.com/c-bata/go-prompt v0.2.5
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stromland/cobra-prompt v0.0.0-20181123224253-940a0a2bd0d3
)
