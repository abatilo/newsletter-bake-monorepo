variable "GITHUB_SHA" {
  default = "latest"
}

variable "REGISTRY" {
  default = "ghcr.io/abatilo/newsletter-bake-monorepo"
}

group "default" {
  targets = [
    "app1",
    "app2",
  ]
}

target "node-modules" {
  target = "node-modules"
  cache-from = ["type=gha,scope=node-modules"]
  cache-to = ["type=gha,mode=max,scope=node-modules"]
}

target "go-modules" {
  target = "go-modules"
  cache-from = ["type=gha,scope=go-modules"]
  cache-to = ["type=gha,mode=max,scope=go-modules"]
}

target "app1" {
  contexts = {
    node-modules = "target:node-modules",
    go-modules = "target:go-modules"
  }
  args = {
    app = "app1"
  }
  tags = [
    "${REGISTRY}/app1:${GITHUB_SHA}",
  ]
  cache-from = ["type=gha,scope=app1"]
  cache-to = ["type=gha,mode=max,scope=app1"]
}

target "app2" {
  contexts = {
    node-modules = "target:node-modules",
    go-modules = "target:go-modules"
  }
  args = {
    app = "app2"
  }
  tags = [
    "${REGISTRY}/app2:${GITHUB_SHA}",
  ]
  cache-from = ["type=gha,scope=app2"]
  cache-to = ["type=gha,mode=max,scope=app2"]
}
