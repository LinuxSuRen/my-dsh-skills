// Package skills embeds the skill directories of this repository and
// provides APIs to enumerate and install them for AI coding tools.
package skills

import (
	"embed"
	"io/fs"
)

// Skill directories at the repository root. A new skill must be added
// to this list to get embedded into the binary.
//
//go:embed android-camera-dev android-device-debugging arch-dev-workflow async-job-sse-upload-dev branching-workflow edge-offline-pipeline-dev frontend-dev-workflow git-identity mqtt-device-gateway-dev ohos-media-dev onvif-rtsp-dev open-source-contribution skill-distillation skill-refresh
var fsys embed.FS

// FS returns the read-only filesystem with all embedded skills.
func FS() fs.FS { return fsys }
