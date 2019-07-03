#!/usr/bin/env groovy
def PROJECT_VERSION = "v0.0.${BUILD_NUMBER}"

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
          bat "buildInJenkins.stage3.cmd ${PROJECT_VERSION}"
        }
      }
    }
  }
}
