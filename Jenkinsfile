#!/usr/bin/env groovy
properties([
  buildDiscarder(logRotator(numToKeepStr: '3')),
  disableConcurrentBuilds()]
)

def PROJECT_VERSION = "v0.0.${BUILD_NUMBER}"

node('che-windows-02') {
  try {
      subst{
        stage ('build') {
          dir('scripts') {
            echo "Running ${env.BUILD_ID} on ${env.JENKINS_URL} ${env.GITHUB_TOKEN}"
            bat 'buildInJenkins.stage1.cmd'
            bat 'buildInJenkins.stage2.cmd'
            bat "buildInJenkins.stage3.cmd ${PROJECT_VERSION} ${env.GITHUB_TOKEN}"
        }
      }
    }
  }
  catch (Exception e) {
    echo 'Exception thrown'
    throw e
  }
}