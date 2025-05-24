#!/bin/bash
set -eo
export AWS_PROFILE=reetro_application
aws ecr get-login-password --region ap-south-1 | docker login --username USERNAME --password-stdin PASSWORD

cd reetro-application

git checkout main
git pull origin main

make prod-up
