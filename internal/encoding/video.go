package encoding

import (
	"cadenceforge/internal/domain"
	"fmt"
)

// videoChain builds the per-placement video normalization used by MP4
// renders: trim and PTS reset matching the audio span, then uniform framing,
// aspect-preserving scale with padding, SAR, frame rate and pixel format.
func videoChain(i int, p domain.Placement, r domain.Recipe) string {
	return fmt.Sprintf("[%d:v:0]trim=start=%s:end=%s,setpts=PTS-STARTPTS,scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1,fps=%d,format=yuv420p[v%d]",
		i, seconds(p.SourceInMS), seconds(p.SourceOutMS), r.Width, r.Height, r.Width, r.Height, r.FPS, i)
}
