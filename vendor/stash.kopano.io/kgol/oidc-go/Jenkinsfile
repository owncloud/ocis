#!/usr/bin/env groovy

pipeline {
	agent {
		docker {
			image 'golang:1.13'
		}
	}
	environment {
		GOBIN = '/tmp/go-bin'
		GOCACHE = '/tmp/go-build'
	}
	stages {
		stage('Bootstrap') {
			steps {
				echo 'Bootstrapping..'
				sh 'go version'
				sh 'go get -v golang.org/x/lint/golint'
				sh 'go get -v github.com/tebeka/go2xunit'
				sh 'go get -v github.com/axw/gocov/...'
				sh 'go get -v github.com/AlekSi/gocov-xml'
				sh 'go mod vendor'
			}
		}
		stage('Lint') {
			steps {
				echo 'Linting..'
				sh 'PATH=$PATH:$GOBIN golint | tee golint.txt || true'
				sh 'go vet | tee govet.txt || true'
			}
		}
		stage('Test') {
			steps {
				echo 'Testing..'
				sh 'PATH=$PATH:$GOBIN go test -v -count=1 -covermode=atomic -coverprofile=coverage.out | tee tests.output'
				sh 'PATH=$PATH:$GOBIN go2xunit -fail -input tests.output -output tests.xml'
			}
		}
		stage('Coverage') {
			steps {
				echo 'Coverage..'
				sh 'mkdir -p ./test/reports'
				sh 'go tool cover -html=coverage.out -o test/reports/coverage.html'
				sh 'PATH=$PATH:$GOBIN; gocov convert coverage.out | gocov-xml > coverage.xml'
				publishHTML([allowMissing: true, alwaysLinkToLastBuild: true, keepAll: true, reportDir: 'test/reports', reportFiles: 'coverage.html', reportName: 'Go Coverage Report HTML', reportTitles: ''])
				step([$class: 'CoberturaPublisher', autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'coverage.xml', failUnhealthy: false, failUnstable: false, maxNumberOfBuilds: 0, onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false])
			}
		}
	}
	post {
		always {
			junit allowEmptyResults: true, testResults: 'tests.xml'
			recordIssues qualityGates: [[threshold: 100, type: 'TOTAL', unstable: true]], tools: [goVet(pattern: 'govet.txt'), goLint(pattern: 'golint.txt')]
			cleanWs()
		}
	}
}
