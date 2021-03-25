#!/bin/bash
#sh copy_data.sh 2020-11-30 20201130
mkdir $1_data
tar -xzvf /home/yangkan/rxy/data/$1.tar.gz
cp $1/container_memory_usage_bytes.log $1_data/
cp $1/container_tasks_state.log $1_data/
cp $1/dcgm_fb_used.log $1_data/
cp $1/dcgm_fb_free.log $1_data/
cp $1/dcgm_gpu_utilization.log $1_data/
cp $1/kube_pod_container_resource_limits.log $1_data/
#cp $1/rate\(container_cpu_usage_seconds_total%5B1m%5D\).log $1_data/
cp $1/rate* $1_data/
#date="PodLifecycle_log."$2"0000"
#cp /home/yangkan/rxy/data/PodLifecycle_log/$date $1_data/PodLifecycle_log.log
rm -rf $1
mv $1_data $1
