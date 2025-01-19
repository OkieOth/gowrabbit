module examples/basic

go 1.23.1

require (
	github.com/okieoth/gowrabbit/pub v0.0.0
	github.com/okieoth/gowrabbit/sub v0.0.0
)

replace github.com/okieoth/gowrabbit/sub => ../../sub

replace github.com/okieoth/gowrabbit/pub => ../../pub
