#!/bin/bash
set -eo
export AWS_PROFILE=reetro_golang
aws ecr get-login-password --region ap-south-1 | docker login --username AWS --password-stdin 576641250044.dkr.ecr.ap-south-1.amazonaws.com

cd reetro-golang

git checkout main
git pull origin main

make prod-up
