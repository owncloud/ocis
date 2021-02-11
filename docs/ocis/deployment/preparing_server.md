---
title: "Preparing a server"
date: 2020-10-12T14:04:00+01:00
weight: 100
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: preparing_server.md
---

{{< toc >}}


## Example for Hetzner Cloud
* create server on Hetzner Cloud. Set labels "owner" and "for". Example for hcloud cli:
`hcloud server create --type cx21 --image ubuntu-20.04 --ssh-key admin --name ocis-server --label owner=admin --label for=testing` 

* Configure DNS A-records for needed domains pointing on the servers ip address, for example in CloudFlare

* Access server via ssh as root

* Create a new user

  `$ adduser --disabled-password --gecos "" admin`

* Add user to sudo group

  `$ usermod -aG sudo admin`

* Install docker

  ```
  apt update
  apt install docker.io
  ```

* Add user to docker group

  `usermod -aG docker admin`

* Install docker-compose via

  `curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose`

  (docker compose version 1.27.4 as of today)
* Make docker-compose executable

  `chmod +x /usr/local/bin/docker-compose`


* Add users pub key to 
    ```
    mkdir /home/admin/.ssh
    echo "<pubkey>" >> /home/admin/.ssh/authorized_keys`
    chown admin:admin -R /home/admin/.ssh
    ```

* Secure ssh daemon by editing `/etc/ssh/sshd_config`
    ```
    PermitRootLogin no
    ChallengeResponseAuthentication no
    PasswordAuthentication no
    UsePAM no
    ```

* restart sshd server to apply settings `systemctl restart sshd`

* Login as the user you created
