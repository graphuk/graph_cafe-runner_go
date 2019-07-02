#!/usr/bin/env groovy
def PROJECT_VERSION = "0.0.${BUILD_NUMBER}"

pipeline {
  agent { label 'windows-dockerized' }
  options {
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '5'))
  }
  stages {
    when {
      anyOf {
        branch 'master';
      }
    }
    stage('BuildAndTest') {
      steps {
        bat 'dir'
      }
    }
    stage('Test') {
      steps {
        bat 'node --version'
      }
    }
  }
}
