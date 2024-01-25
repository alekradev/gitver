#!/bin/bash

go build -o build/gitver
sudo cp build/gitver /usr/bin/gitver
