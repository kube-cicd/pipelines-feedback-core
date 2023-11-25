package k8s

import (
	"bytes"
	"context"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"io"
	v1api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"sort"
	"strconv"
	"strings"
)

// FindAndReadLogsFromLastPod is listing logs from the last Pod found for given selector. Errors are returned as logs
// 1. Find all Pods for given Job
// 2. Sort by metadata.creationTimestamp, desc
// 3. Pick first and retrieve logs
func FindAndReadLogsFromLastPod(ctx context.Context, lister v1.PodInterface, selector string) string {
	podList, err := lister.List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return "Pipelines Feedback Core: Cannot list Pods for selector: " + err.Error()
	}

	listAsArr := podList.Items
	sort.Slice(listAsArr, func(x, y int) bool {
		return listAsArr[y].CreationTimestamp.Before(&listAsArr[x].CreationTimestamp)
	})

	if len(listAsArr) == 0 {
		return "Pipelines Feedback Core: No Pods found for selector"
	}

	req := lister.GetLogs(listAsArr[0].Name, &v1api.PodLogOptions{})
	return ReadRequestStream(ctx, req)
}

// ReadRequestStream is a helper you can use to read logs from the Pod. Errors are returned as logs
func ReadRequestStream(ctx context.Context, req *rest.Request) string {
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "Pipelines Feedback Core: Cannot open stream: " + err.Error()
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "Pipelines Feedback Core: Cannot copy stream: " + err.Error()
	}
	return buf.String()
}

// TruncateLogs is truncating logs with a maximum lines number, maximum line length
func TruncateLogs(logs string, data config.Data) string {
	maxLineLength, _ := strconv.Atoi(data.GetOrDefault("logs-max-line-length", "64"))
	maxFullLengthLines, _ := strconv.Atoi(data.GetOrDefault("logs-max-full-length-lines-count", "10"))
	lineSplitSeparator := data.GetOrDefault("logs-split-separator", "(...)")
	maxLogsLength := (maxFullLengthLines * maxLineLength) + (maxFullLengthLines * len(lineSplitSeparator))

	asLines := strings.Split(logs, "\n")
	lines := make([]string, 0)
	startingFrom := len(asLines) - maxFullLengthLines
	if startingFrom < 0 {
		startingFrom = 0
	}

	processed := 0
	for num, line := range asLines {
		if num < startingFrom {
			continue
		}
		processed += 1
		if len(line) > maxLineLength+len(lineSplitSeparator) {
			firstPartEnds := maxLineLength / 2
			secondPartStarts := len(line) - (maxLineLength / 2)
			lines = append(lines, line[0:firstPartEnds]+lineSplitSeparator+line[secondPartStarts:])
		} else {
			lines = append(lines, line)
		}
		if processed > maxFullLengthLines {
			break
		}
	}
	logs = strings.Join(lines, "\n")
	if len(logs) < maxLogsLength {
		return logs
	}
	return logs[len(logs)-maxLogsLength:]
}
