// Package build_info
//
// ビルド時に動的に埋め込まれる設定値を管理するパッケージです。
package build_info

var (
	serverVersion string
)

func ServerVersion() string {
	return serverVersion
}
