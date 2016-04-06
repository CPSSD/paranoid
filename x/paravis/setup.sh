#!/bin/bash

# Download the latest version of jQuery
wget -O lib/jquery-2.2.3.min.js https://code.jquery.com/jquery-2.2.3.min.js

# Download required fonts. Please periodically make sure the fonts are available
mkdir -p ./res/fonts
mkdir -p /tmp/paravis
wget -O /tmp/paravis/Roboto.zip https://material-design.storage.googleapis.com/publish/material_v_4/material_ext_publish/0B0J8hsRkk91LRjU4U1NSeXdjd1U/RobotoTTF.zip
unzip /tmp/paravis/Roboto.zip -d /tmp/paravis/roboto-extracted
cp /tmp/paravis/roboto-extracted/Roboto-Regular.ttf ./res/fonts/
cp /tmp/paravis/roboto-extracted/Roboto-Medium.ttf ./res/fonts/
cp /tmp/paravis/roboto-extracted/Roboto-Thin.ttf ./res/fonts/
rm -r /tmp/paravis
