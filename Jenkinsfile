#!/usr/bin/env groovy
def PROJECT_VERSION = "0.0.${BUILD_NUMBER}"

pipeline {
  agent { label 'windows-dockerized' }
  options {
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '5'))
  }
  stages {
    stage('BuildAndTest') {
      when {
        anyOf {
          branch 'master';
        }
      }
      steps {
        dir('scripts') {
          bat 'buildInJenkins.stage1.cmd'
          bat 'buildInJenkins.stage2.cmd'
        }
      }
    }
    stage('Test') {
      when {
        anyOf {
          branch 'master';
        }
      }
      steps {
        bat 'node --version'
      }
    }
  }
}
