#!/bin/groovy
@Library("shared") _

def LABEL = "pod-${UUID.randomUUID().toString()}"
def standsToDeploy

podTemplate(
  label: LABEL,
  inheritFrom: "default"
) {
  node(LABEL) {
    properties([
      parameters([
        stringParam(
          defaultValue: "latest",
          description: "DOCKER_TAG to deploy",
          name: "dockerTag"
        )
      ])
    ])

    try {
      checkout scm

      stage("Ask") {
        standsToDeploy = askStand(5, "MINUTES", "STANDS", ["dev", "stage", "prod"])
      }

      if (standsToDeploy.isEmpty()) {
        println("Не выбрал куда деплоить :/")
        currentBuild.result = "SUCCESS"
        return
      }

      stage("Deploy") {
        container("kubectl-helm3-default") {
          standsToDeploy.each { envName ->
            deploy(envName) {
              dir("./deploy") {
                println("deploy to ${envName}")
                sh "helm upgrade --install maas-api ./maasapi --namespace ${envName} -f values-${envName}.yaml --set image.tag=${params.dockerTag}"
              }
            }
          }
        }
      }

      currentBuild.result = "SUCCESS"
    } catch (e) {
      println(e)
      currentBuild.result = "FAILURE"
    }
  }
}