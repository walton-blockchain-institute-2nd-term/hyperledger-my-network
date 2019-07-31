#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -e

# Shut down the Docker containers for the system tests.
docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

# remove the local state
rm -f ~/.hfc-key-store/*

# remove chaincode docker images
# f 옵션은 파일
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images dev-* -q)

sleep 1
# 사용하지 않는 네트워크 삭제
docker network prune
# Your system is now clean
