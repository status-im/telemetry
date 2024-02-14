pipeline {
  agent { label 'linux' }

  parameters {
    booleanParam(
      name: 'DEPLOY',
      description: 'Enable to deploye the Docker image.',
      defaultValue: false,
    )
    string(
      name: 'DOCKER_CRED',
      description: 'Name of Docker Registry credential.',
      defaultValue: params.DOCKER_CRED ?: 'harbor-telemetry-robot',
    )
    string(
      name: 'DOCKER_REGISTRY_URL',
      description: 'URL of the Docker Registry',
      defaultValue: params.DOCKER_REGISTRY_URL ?: 'https://harbor.status.im'
    )
  }

  options {
    timestamps()
    disableConcurrentBuilds()
    /* manage how many builds we keep */
    buildDiscarder(logRotator(
      numToKeepStr: '10',
      daysToKeepStr: '30',
    ))
  }

  environment {
    IMAGE_NAME = "status-im/telemetry"
    IMAGE_DEFAULT_TAG = "${env.GIT_COMMIT.take(7)}"
    IMAGE_DEPLOY_TAG = "deploy"
  }

  stages {
    stage('Build') { steps { script {
      image = docker.build(
        "${env.IMAGE_NAME}:${env.IMAGE_DEFAULT_TAG}"
      )
    } } }

    stage('Push') { steps { script {
      withDockerRegistry([
        credentialsId: params.DOCKER_CRED, url: params.DOCKER_REGISTRY_URL
      ]) {
        image.push()
      }
    } } }

    stage('Deploy') {
      when { expression { params.DEPLOY } }
      steps { script {
        withDockerRegistry([
          credentialsId: params.DOCKER_CRED, url: params.DOCKER_REGISTRY_URL
        ]) {
          image.push(env.IMAGE_DEPLOY_TAG)
        }
    } } }
  } // stages
}
