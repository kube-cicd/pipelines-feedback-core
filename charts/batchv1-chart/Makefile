build:
	cd .. && rsync -av ./pipelines-feedback-common/* ./batchv1-chart/ --exclude Makefile
	cd ../batchv1-chart && yq e -i '.name= "batchv1-chart"' Chart.yaml
