package errors

const (
	AnnotationCannotParseURL ErrorType = -(1800 + iota)
	AnnotationCannotGetFromServer
	AnnotationCannotReadBody
)
