package encoding

// encoderArgs selects stream mapping and codec flags for the output container.
// It only decides how the already-built [aout] (and [vout]) streams are
// packaged; it contains no filter graph or loudness semantics.
func encoderArgs(format string) []string {
	if format == "mp4" {
		return []string{
			"-map", "[vout]", "-map", "[aout]",
			"-c:v", "libx264", "-preset", "veryfast", "-crf", "23", "-threads", "1",
			"-c:a", "aac", "-b:a", "160k", "-movflags", "+faststart",
		}
	}
	return []string{"-map", "[aout]", "-c:a", "pcm_s16le"}
}
