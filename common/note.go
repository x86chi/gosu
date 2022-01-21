package common

// time rate has been applied in advance
type Note struct {
	Type  NoteType
	Time  int64
	Time2 int64 // ex) ln end time

	SampleVolume   int
	SampleFilename string
}

type NoteType int16