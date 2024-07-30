#!/bin/bash

#SBATCH --job-name=$2
##SBATCH --output=myjob.out
##SBATCH --error=myjob.err
#SBATCH --partition=testbatch
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1

echo $1
echo "Hello, world!"

sleep 30
