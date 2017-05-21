pipeline {
    agent {
        docker {
            image 'golang:1.8'
        }
    }
    environment { 
        GIT_COMMITTER_NAME = 'jenkins'
        GIT_COMMITTER_EMAIL = 'jenkins@localhost'
        GIT_AUTHOR_NAME = 'jenkins'
        GIT_AUTHOR_EMAIL = 'jenkins@localhost'
        GOPATH = "${WORKSPACE}"
    }
    stages {
        stage('Test') {
            steps {
                sh '''
                    cd ${WORKSPACE}
                    go get github.com/golang/lint/golint
                    go get github.com/tebeka/go2xunit
                    go get git.yale.edu/docker/deco/...
                    ./bin/golint src/git.yale.edu/docker/deco/.. > lint.txt
                    go test -v git.yale.edu/docker/deco/... > test.output
                    # fails with no tests, re-enable when tests are added
                    # sh 'cat test.output | ./bin/go2xunit -output tests.xml'
                '''
            }
            post {
                success {
                    stash includes: 'lint.txt,test.output,test.xml', name: 'decoTests'
                }
            }
        }
        stage('Build'){
            steps {
                sh 'cd $WORKSPACE && go build -o deco-native -v git.yale.edu/docker/deco'
                sh './deco-native version -s > deco.version'
                sh '''
                    cd ${WORKSPACE}
                    VERSION=`cat deco.version`
                    [[ !  -z  ${VERSION}  ]] && echo 'VERSION not found' && exit 1
                    for GOOS in darwin linux; do
                        for GOARCH in 386 amd64; do
                            echo "Building $GOOS-$GOARCH"
                            export GOOS=$GOOS
                            export GOARCH=$GOARCH
                            go build -o deco-v${VERSION}-$GOOS-$GOARCH git.yale.edu/docker/deco
                        done
                    done
                '''
            }
            post {
                success {
                    stash includes: 'deco*', name: 'decoBin'
                }
            }
        }
    }
    options {
        buildDiscarder(logRotator(numToKeepStr:'3'))
        timeout(time: 60, unit: 'MINUTES')
    }
}
