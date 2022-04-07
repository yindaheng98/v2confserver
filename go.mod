module github.com/yindaheng98/v2confserver

go 1.16

replace github.com/v2fly/vmessping v0.3.5-0.20211004134616-eb37e6100b2a => github.com/yindaheng98/vmessping v0.3.5-0.20220405100036-6c27aab8ef0a

require github.com/yindaheng98/vmessconfig v0.0.0-20220407044324-0c4cb5b3f925 // indirect

require (
	github.com/octago/sflags v0.3.1-0.20210726012706-20f2a9c31dfc // indirect
	github.com/v2fly/v2ray-core/v4 v4.43.0
)
