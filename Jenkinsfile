pipeline {
    agent { label 'agent1' }
    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'pwd'
                sh 'go mod download && go mod verify'
                sh 'go build -v -o app ./...'
                archiveArtifacts artifacts: 'app'
            }
        }
    }
}
