# Local mt protocol snapshot

This directory starts from `github.com/HimbeerserverDE/mt` commit
`7d1d3af7c7d4` (the proxy's Luanti 5.17 dependency). The upstream module has
not yet added Luanti 5.18's extra `TOCLIENT_SET_LIGHTING` fields.

The Skyss changes are in `tocltcmds.go` and the corresponding generated
`serialize.go`: a presence byte, an optional 16-byte light curve, and a
32-bit ambient occlusion gamma value. `lighting_518_test.go` checks both wire
shapes. Keep this module and the proxy's `go.mod` replacement together until
upstream `mt` supports protocol 54.
