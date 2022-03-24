build:
	docker build -t jefedavis/resume-pdf-converter:latest .

push:
	docker push jefedavis/resume-pdf-converter:latest