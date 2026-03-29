
def SERVICE_NAME = 'rest-api'
def imageTag = ""

pipeline {
    // Changed from 'kubernetes' to 'any' to use the local Jenkins node environment
    agent any 

    //  tools {
        
    //          go 'go' 

    //     }

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

        stage('Determine Version') {
            steps {
                
                    script {
                        sh 'git config --global --add safe.directory "$(pwd)"'

                        def branchName = env.BRANCH_NAME
                        def gitCommit = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                        def isExactTag = sh(script: 'git describe --exact-match --tags HEAD 2>/dev/null || true', returnStdout: true).trim()

                        echo "Branch: ${branchName}"
                        echo "Commit: ${gitCommit}"
                        echo "Exact tag: ${isExactTag ?: 'none'}"

                        if (isExactTag) {
                            // Production release build
                            imageTag = isExactTag
                        } else {
                            // Dev / staging / feature builds
                            def safeBranch = branchName
                                .replace('/', '-')
                                .replace('_', '-')
                                .toLowerCase()

                            imageTag = "${safeBranch}-${gitCommit}"
                        }

                        echo "Using image tag: ${imageTag}"

                        env.IMAGE_TAG = imageTag
                        env.GIT_COMMIT = gitCommit
                    }
                
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
                    sh "sleep 10"
                
                    env.SONAR_HOST_URL = "https://code.touchvas.work" 
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