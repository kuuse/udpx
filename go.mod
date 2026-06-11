module github.com/nullt3r/udpx

go 1.17

require (
	github.com/nullt3r/udpx/pkg/colors v0.0.0
	github.com/nullt3r/udpx/pkg/probes v0.0.0
	github.com/nullt3r/udpx/pkg/scan v0.0.0
	github.com/nullt3r/udpx/pkg/targets v0.0.0
	github.com/nullt3r/udpx/pkg/utils v0.0.0
	github.com/nullt3r/udpx/pkg/version v0.0.0
)

replace (
	github.com/nullt3r/udpx/pkg/colors => ./pkg/colors
	github.com/nullt3r/udpx/pkg/probes => ./pkg/probes
	github.com/nullt3r/udpx/pkg/scan => ./pkg/scan
	github.com/nullt3r/udpx/pkg/targets => ./pkg/targets
	github.com/nullt3r/udpx/pkg/utils => ./pkg/utils
	github.com/nullt3r/udpx/pkg/version => ./pkg/version
)
