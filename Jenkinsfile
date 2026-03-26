def NEXUS_DOCKER_CREDENTIALS = 'nexus-creds'
def NEXUS_DOCKER_REPO = 'registry.touchvas.work'
def NEXUS_DOCKER_REPO_URL = 'https://registry.touchvas.work'
def SERVICE_NAME = 'rest-api'
def imageTag = ""

pipeline {
    // Changed from 'kubernetes' to 'any' to use the local Jenkins node environment
    agent any 

     tools {
        
             go 'go' 

        }

    stages {
        stage('Checkout Code') {
            steps { checkout scm }
        }

    

        stage('Setup Dependencies') {
            steps {
            
                sh 'go mod tidy'
                sh 'go mod vendor'
            }
        }

        stage('Unit Tests') {
            steps {
                sh 'go test'
            }
        }

        stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'SonarQubeScanner'
                    withSonarQubeEnv('SonarQube') {
                        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${SERVICE_NAME}"
                    }
                }
            }
        }

        stage("Quality Gate Check") {
            steps {
                script {
                    sh "sleep 5"
                
                    env.SONAR_HOST_URL = "http://sonarqube:9000" 
                }
                timeout(time: 15, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Build and Push Docker Image to Nexus') {

            when {
                 expression { return false } 
            }

            steps {
                script {
                    // Force a branch name if testing locally where env.BRANCH_NAME might be null
                    def branchName = env.BRANCH_NAME ?: "dev" 

                    if (branchName == 'master') {
                        imageTag = 'latest'
                    } else if (branchName == 'dev') {
                        imageTag = 'dev'
                    } else {
                        echo "Skipping pushing branch: ${branchName}"
                        return
                    }

                    def fullImageName = "${NEXUS_DOCKER_REPO}/touchvas/${SERVICE_NAME}:${imageTag}"

                    
                    withCredentials([usernamePassword(
                        credentialsId: NEXUS_DOCKER_CREDENTIALS,
                        usernameVariable: 'NEXUS_USER',
                        passwordVariable: 'NEXUS_PASS')]) 
                    {
                        // Removed container('docker-cli') wrapper
                        sh "echo \"${NEXUS_PASS}\" | docker login -u ${NEXUS_USER} --password-stdin ${NEXUS_DOCKER_REPO_URL}"
                        sh "docker build --no-cache -t ${fullImageName} ."
                        sh "docker push ${fullImageName}"
                        sh "docker logout ${NEXUS_DOCKER_REPO_URL}"
                    }
                }
            }
        }

        // Optional: Keep or Comment out depending on if your local machine can reach the K8s cluster
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    echo "Skipping K8s deployment for local test"
                }
            }
        }
    }
}