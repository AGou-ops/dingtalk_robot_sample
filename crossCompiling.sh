#!/bin/bash
#
#**************************************************
# Author:         AGou-ops                        *
# E-mail:         agou-ops@foxmail.com            *
# Date:           2021-12-09                      *
# Description:                              *
# Copyright 2021 by AGou-ops.All Rights Reserved  *
#**************************************************

binaryBaseName="DingTalk_robot"
version="v0.11"
dist_archs=(386 amd64 arm arm64 mips mips64 mips64le mipsle ppc64 ppc64le riscv64 s390x)

rm -rf ../output_binary
mkdir ../output_binary

for arch in ${dist_archs[@]}
do
        env GOOS=linux GOARCH=${arch} go build -x -v -o ../output_binary/${binaryBaseName}_${version}_linux_${arch}
done

env GOOS=windows GOARCH=amd64 go build -x -v -o ../output_binary/${binaryBaseName}_${version}_windows_amd64

trap "echo 'program exit...'; exit 2" SIGINT

echo -e "\n\n"

echo "生成checksum..."

(cd ../output_binary && shasum * > ${binaryBaseName}_${version}.checksums.txt)


echo "
===========
=         =
=  Done.  =
=         =
===========
"
