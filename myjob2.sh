#!/bin/bash

#SBATCH --job-name=mailtest
##SBATCH --output=myjob.out
##SBATCH --error=myjob.err
#SBATCH --partition=batch
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --mail-user=jyjungs63@gmail.com
#SBATCH --mail-type=BEGIN
#SBATCH --mail-type=END
#SBATCH --mail-type=FAIL

echo $1
echo "Hello, world!"

sleep 30
