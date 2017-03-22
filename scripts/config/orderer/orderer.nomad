job "orderer" {
  datacenters = ["dc1"]
  region = "global"
  type = "service"
  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }
  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "orderers" {
    count = 1

    restart {
      attempts = 5
      interval = "30s"
      delay = "5s"
      mode = "delay"
    }

    ephemeral_disk {
      size = 300
    }

    task "orderer" {
      driver = "docker"
      
      config {
        image = "ncodes/cocoon-launcher:latest"
        command = "bash"
        args = ["run.sh"]
        work_dir = "/local/scripts"
        port_map {}
      }

      artifact {
        source = "https://raw.githubusercontent.com/ncodes/cocoon/master/scripts/config/orderer/run.sh"
        destination = "/local/scripts"
      }

      logs {
        max_files     = 10
        max_file_size = 10
      }

      resources {
        cpu    = 500
        memory = 256 
        network {
          mbits = 1000
          port "orderer_grpc" {}
        }
      }

      service {
        name = "orderer"
        tags = []
        port = "orderer_grpc"
        check {
          name     = "alive"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}