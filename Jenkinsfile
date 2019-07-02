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
        bat 'mkdir src\\github.com\\graph-uk'
        bat 'mklink /D src\\github.com\\graph-uk\\graph_cafe-runner_go %CD%'
        dir('src\\github.com\\graph-uk\\graph_cafe-runner_go') {
          bat 'npm install'
          bat 'buildReleaseAndTestIntegration.cmd'
        }
        bat 'rd src\\github.com\\graph-uk\\graph_cafe-runner_go'
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
