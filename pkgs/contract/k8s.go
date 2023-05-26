package contract

type PipelineAnnotation string

const (
	AnnotationPrId         = "pipelinesfeedback.keskad.pl/pr-id"
	AnnotationCommitHash   = "pipelinesfeedback.keskad.pl/commit"
	AnnotationHttpsRepo    = "pipelinesfeedback.keskad.pl/https-repo-url"
	AnnotationReference    = "pipelinesfeedback.keskad.pl/ref"
	AnnotationTechnicalJob = "pipelinesfeedback.keskad.pl/technical-job"
)
