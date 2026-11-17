package encoding

import (
	"cadenceforge/internal/domain"
	"cadenceforge/internal/media"
	"fmt"
	"strings"
)

func seconds(ms int64) string { return fmt.Sprintf("%.3f", float64(ms)/1000) }

// concatChain joins the per-placement labelled streams in timeline order.
// Audio labels are joined for every render; the video labels are included
// only when the MP4 container carries a picture. The concat filter variant is
// selected from the container, but the audio labels themselves are produced
// by audioChain without any container knowledge.
func concatChain(format string, count int, joins []string) string {
	if format == "mp4" {
		return strings.Join(joins, "") + fmt.Sprintf("concat=n=%d:v=1:a=1[vout][acat]", count)
	}
	return strings.Join(joins, "") + fmt.Sprintf("concat=n=%d:v=0:a=1[acat]", count)
}

// Arguments assembles the full ffmpeg command. Filter semantics and container
// selection are separate responsibilities:
//   - audioChain / loudnessChain own every audio filter and normalization;
//   - videoChain owns MP4 picture normalization;
//   - concatChain wires the labelled streams;
//   - encoderArgs only maps the finished labels onto container codecs.
//
// Clip timing, edit gain and the loudness policy are therefore independent of
// whether the result is muxed as MP4 or WAV.
func Arguments(render domain.Render, paths map[string]string, output string) ([]string, error) {
	r := render.Recipe
	args := []string{"-nostdin", "-y", "-hide_banner", "-loglevel", "error", "-filter_complex_threads", "1"}
	filters := []string{}
	joins := []string{}
	for i, p := range render.Plan.Placements {
		path, ok := paths[p.SourceID]
		if !ok {
			return nil, fmt.Errorf("source object missing")
		}
		args = append(args, "-protocol_whitelist", "file,pipe", "-format_whitelist", media.Formats, "-i", path)
		filters = append(filters, audioChain(i, p, r.SampleRate, render.Sources[p.SourceID].Media.Audio))
		if r.Format == "mp4" {
			filters = append(filters, videoChain(i, p, r))
			joins = append(joins, fmt.Sprintf("[v%d][a%d]", i, i))
		} else {
			joins = append(joins, fmt.Sprintf("[a%d]", i))
		}
	}
	if len(joins) == 0 {
		return nil, fmt.Errorf("empty render")
	}
	filters = append(filters, concatChain(r.Format, len(joins), joins))
	filters = append(filters, loudnessChain(r))
	args = append(args, "-filter_complex", strings.Join(filters, ";"))
	args = append(args, encoderArgs(r.Format)...)
	return append(args, "-ar", fmt.Sprint(r.SampleRate), "-t", seconds(render.Plan.DurationMS), "-map_metadata", "-1", output), nil
}
