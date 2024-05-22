package contract

import (
	"os"
	"strings"
)

// GetPrIdAnnotation returns by default "pipelinesfeedback.keskad.pl/pr-id". Parametrized with 'ANNOTATION_FEEDBACK_BASE' env variable
func GetPrIdAnnotation() string {
	return getAnnotationBase() + "/pr-id"
}

// GetCommmitAnnotation returns by default "pipelinesfeedback.keskad.pl/commit". Parametrized with 'ANNOTATION_FEEDBACK_BASE' env variable
func GetCommmitAnnotation() string {
	return getAnnotationBase() + "/commit"
}

// GetHttpsRepoUrlAnnotation returns by default "pipelinesfeedback.keskad.pl/https-repo-url". Parametrized with 'ANNOTATION_FEEDBACK_BASE' env variable
func GetHttpsRepoUrlAnnotation() string {
	return getAnnotationBase() + "/https-repo-url"
}

// GetRefAnnotation returns by default "pipelinesfeedback.keskad.pl/ref". Parametrized with 'ANNOTATION_FEEDBACK_BASE' env variable
func GetRefAnnotation() string {
	return getAnnotationBase() + "/ref"
}

// GetTechnicalJobAnnotation returns by default "pipelinesfeedback.keskad.pl/technical-job". Parametrized with 'ANNOTATION_FEEDBACK_BASE' env variable
func GetTechnicalJobAnnotation() string {
	return getAnnotationBase() + "/technical-job"
}

func getAnnotationBase() string {
	if val := os.Getenv("ANNOTATION_FEEDBACK_BASE"); val != "" {
		return os.Getenv("ANNOTATION_FEEDBACK_BASE")
	}
	return "pipelinesfeedback.keskad.pl"
}

func getFeedbackLabel() (string, string) {
	labelName := "pipelinesfeedback.keskad.pl/enabled"
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
	requiredLabelName, requiredLabelValue := getFeedbackLabel()
	// there is a label present
	if val, ok := labels[requiredLabelName]; ok {
		// label has required value set
		return strings.Trim(strings.ToLower(val), " ") == requiredLabelValue
	}
	// no label present
	return false
}
