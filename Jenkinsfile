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
        bat 'IF EXIST src rd /s /q src'
        bat 'mkdir src\\github.com\\graph-uk'
        bat 'mklink /D src\\github.com\\graph-uk\\graph_cafe-runner_go %CD%'
        dir('src\\github.com\\graph-uk\\graph_cafe-runner_go') {
          subst{
            echo 'npm install'
            echo 'buildReleaseAndTestIntegration.cmd'
          }
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
