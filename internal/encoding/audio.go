package encoding

import (
	"cadenceforge/internal/domain"
	"fmt"
)

// audioChain builds the audio filter for one placement. It owns the clip's
// sample timing (trim plus PTS reset), the stereo resample contract and the
// edit gain from the sheet. A source without an audio track is replaced by
// silence spanning the same output interval so concat timing is preserved.
// The chain never branches on the output container: its inputs are only the
// placement, the recipe sample rate and the source's stream facts.
func audioChain(i int, p domain.Placement, sampleRate int, hasAudio bool) string {
	if hasAudio {
		return fmt.Sprintf("[%d:a:0]atrim=start=%s:end=%s,asetpts=PTS-STARTPTS,aresample=%d,aformat=sample_fmts=fltp:channel_layouts=stereo,volume=%.3fdB[a%d]",
			i, seconds(p.SourceInMS), seconds(p.SourceOutMS), sampleRate, p.GainDB, i)
	}
	return fmt.Sprintf("anullsrc=r=%d:cl=stereo,atrim=duration=%s,asetpts=PTS-STARTPTS[a%d]",
		sampleRate, seconds(p.OutputOutMS-p.OutputInMS), i)
}

// loudnessChain builds the single-pass loudness stage applied to the
// concatenated audio. Normalization policy comes solely from the recipe, so
// the same audio semantics apply regardless of the chosen container.
func loudnessChain(r domain.Recipe) string {
	if r.Normalize {
		return fmt.Sprintf("[acat]loudnorm=I=%.2f:TP=-1.5:LRA=11,aresample=%d[aout]", r.TargetLUFS, r.SampleRate)
	}
	return "[acat]anull[aout]"
}
