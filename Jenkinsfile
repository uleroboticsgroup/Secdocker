pipeline {
    agent any
    
    tools {
        go 'go-1.11'
    }
    
    environment {
        GO111MODULE = 'on'
    }
    
    options {
        gitLabConnection('roboticsgroup-gitlab')
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git([
                    url: 'https://niebla.unileon.es/DavidFerng/secdocker.git',
                    credentialsId: 'jenkins-gitlab'
                ])
            }
        }
        stage('Build') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'go test -v ./tests'
            }
        }
        stage('SonarCloud') {
          environment {
            SCANNER_HOME = tool 'Sonarqube'
            PROJECT_NAME = "secdocker"
            ORGANIZATION = "default-organization"
          }
          steps {
            withSonarQubeEnv('roboticsgroup-sonarqube') {
                sh '''$SCANNER_HOME/bin/sonar-scanner -Dsonar.organization=$ORGANIZATION \
                -Dsonar.java.binaries=build/classes/java/ \
                -Dsonar.projectKey=$PROJECT_NAME \
                -Dsonar.sources=.'''
            }
          }
        }
  
    }
    
    post {
        always {
            deleteDir()
        }
    }
}