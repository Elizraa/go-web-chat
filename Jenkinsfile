pipeline {
    agent {
        docker {
            image 'golang:1.20.5-alpine'
            args '-p 3000:3000' 
        }
    }

    stages {
        stage('Build') {
            steps {
                sh 'go build'
            }
        }
    }
}
