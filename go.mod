module github.com/kuuse/udpx

go 1.17

// Allow imports using the original upstream path
// All code imports github.com/nullt3r/udpx/pkg/...
// which gets redirected to ./pkg/...
replace github.com/nullt3r/udpx => ./
