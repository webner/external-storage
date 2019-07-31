#!/bin/bash
make
tar czf nfs-provisioner-v1.0.3.cat.tar.gz nfs-provisioner tp-free deploy/kubernetes deploy/systemd
