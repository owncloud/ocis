#!/usr/bin/env groovy

// kcc-go

pipeline {
	agent {
		docker {
			image 'golang:1.13'
		}
	}
	environment {
		GOBIN = '/tmp/go-bin'
		GOCACHE = '/tmp/go-build'
		HOME = '/tmp'
	}
	stages {
		stage('Bootstrap') {
			steps {
				echo 'Bootstrapping..'
				sh 'export'
				sh 'go version'
				sh 'curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOBIN) v1.21.0'
				sh 'go get -v github.com/tebeka/go2xunit'
				sh 'go mod vendor'
			}
		}
		stage('Lint') {
			steps {
				echo 'Linting..'
				sh 'PATH=$PATH:$GOBIN golangci-lint run --modules-download-mode vendor --out-format checkstyle --issues-exit-code 0 > tests.lint.xml'
				checkstyle pattern: 'tests.lint.xml', canComputeNew: false, unstableTotalHigh: '100'
			}
		}
		stage('Test') {
			steps {
				withCredentials([usernamePassword(credentialsId: 'TEST_CREDENTIALS', usernameVariable: 'TEST_USERNAME', passwordVariable: 'TEST_PASSWORD'), string(credentialsId: 'KOPANO_SERVER_DEFAULT_URI', variable: 'KOPANO_SERVER_DEFAULT_URI')]) {
					echo 'Testing..'
					sh 'echo Kopano Server URI: \$KOPANO_SERVER_DEFAULT_URI'
					sh 'echo Kopano Server Username: \$TEST_USERNAME'
					sh 'go test -v -count=1 | tee tests.output'
					sh 'PATH=$PATH:$GOBIN  go2xunit -fail -input tests.output -output tests.xml'
				}
				junit allowEmptyResults: true, testResults: 'tests.xml'
			}
		}
	}
	post {
		always {
			cleanWs()
		}
	}
}
