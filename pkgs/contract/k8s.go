package contract

import (
	"os"
	"strings"
)

const (
	LabelFeedbackEnabled   = "pipelinesfeedback.keskad.pl/enabled"
	AnnotationPrId         = "pipelinesfeedback.keskad.pl/pr-id"
	AnnotationCommitHash   = "pipelinesfeedback.keskad.pl/commit"
	AnnotationHttpsRepo    = "pipelinesfeedback.keskad.pl/https-repo-url"
	AnnotationReference    = "pipelinesfeedback.keskad.pl/ref"
	AnnotationTechnicalJob = "pipelinesfeedback.keskad.pl/technical-job"
)

func GetFeedbackLabel() (string, string) {
	labelName := LabelFeedbackEnabled
	labelValue := "true"

	if val := os.Getenv("LABEL_FEEDBACK_ENABLED_NAME"); val != "" {
		labelName = os.Getenv("LABEL_FEEDBACK_ENABLED_NAME")
	}
	if val := os.Getenv("LABEL_FEEDBACK_ENABLED_VALUE"); val != "" {
		labelValue = os.Getenv("LABEL_FEEDBACK_ENABLED_VALUE")
	}
	return labelName, labelValue
}

// IsJobHavingRequiredLabel decides if a controller should take the resource
func IsJobHavingRequiredLabel(labels map[string]string) bool {
	requiredLabelName, requiredLabelValue := GetFeedbackLabel()
	// there is a label present
	if val, ok := labels[requiredLabelName]; ok {
		// label has required value set
		return strings.Trim(strings.ToLower(val), " ") == requiredLabelValue
	}
	// no label present
	return false
}
