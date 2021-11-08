pipeline {
  agent { label 'linux' }

  parameters {
    booleanParam(
      name: 'DEPLOY',
      description: 'Enable to deploye the Docker image.',
      defaultValue: false,
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
    IMAGE_NAME = "statusteam/telemetry"
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
        credentialsId: "dockerhub-statusteam-auto"
      ]) {
        image.push()
      }
    } } }

    stage('Deploy') {
      when { expression { params.DEPLOY } }
      steps { script {
        withDockerRegistry([
          credentialsId: "dockerhub-statusteam-auto"
        ]) {
          image.push(env.IMAGE_DEPLOY_TAG)
        }
    } } }
  } // stages
} // pipeline
