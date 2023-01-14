package contract

type PipelineAnnotation string

const (
	AnnotationPrId       = "pipelines-feedback.keskad.pl/pr-id"
	AnnotationCommitHash = "pipelines-feedback.keskad.pl/commit"
	AnnotationHttpsRepo  = "pipelines-feedback.keskad.pl/https-repo-url"
	AnnotationReference  = "pipelines-feedback.keskad.pl/ref"
)
