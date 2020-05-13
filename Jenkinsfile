#!/bin/groovy
@Library("shared") _
import ru.movista.constants.Charts

def LABEL = "pod-${UUID.randomUUID().toString()}"
def IMAGE_NAME = "maas-api"
def DOCKER_TAG = ""

podTemplate(
  label: LABEL,
  inheritFrom: "default",
  containers: [
    containerTemplate(name: "golang", image: "golang:1.13-alpine", command: "cat", ttyEnabled: true)
  ],
  volumes: [
    persistentVolumeClaim(claimName: 'golang-cache', mountPath: '/go/pkg/mod')
  ]
) {
  node(LABEL) {
    REPO = checkout scm
    DOCKER_TAG = makeDockerTag(env.BRANCH_NAME, REPO.GIT_COMMIT)
    
    try {

      stage("Build") {
        container("golang") {
          sh("CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o maasapi")
        }
      }

      if (isPR()) {
      	currentBuild.result = "SUCCESS"
        return
      }

      if (isDevelop() || isRelease()) {
        stage("Docker") {
          container("docker-default") {
            dockerBuildTagPush("./Dockerfile", ".", IMAGE_NAME, DOCKER_TAG)
          }
        }
      }

      if (isDevelop()) {
        deployv3("dev", Charts.Front.MAASAPI, DOCKER_TAG)
      }

      currentBuild.result = "SUCCESS"
    } catch (e) {
      print(e)
      currentBuild.result = "FAILURE"
    }
  }
}