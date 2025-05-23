#!/bin/bash

OBA_VERSION="2.5.13-otsf"
MYSQL_CONNECTOR_VERSION="8.3.0"
POSTGRESQL_CONNECTOR_VERSION="42.7.2"

# Output directory
OUTPUT_DIR="./oba-artifacts"
mkdir -p "$OUTPUT_DIR"

# Function to download and rename artifact
copy_and_rename_artifact() {
    ARTIFACT=$1
    TARGET_NAME=$2

    # Download the artifact to local repo
    mvn dependency:get -Dartifact="$ARTIFACT"

    # Copy the artifact to the target directory
    mvn dependency:copy -Dartifact="$ARTIFACT" -DoutputDirectory=$OUTPUT_DIR -Dmdep.useBaseVersion=true

    # Rename the file
    FILE_NAME=$(basename "$OUTPUT_DIR"/*"$(echo "$ARTIFACT" | awk -F: '{print $2}')"*)
    EXT="${FILE_NAME##*.}"
    mv "$OUTPUT_DIR/$FILE_NAME" "$OUTPUT_DIR/$TARGET_NAME.$EXT"
}

copy_and_rename_artifact \
    "org.onebusaway:onebusaway-api-webapp:${OBA_VERSION}:war" \
    "onebusaway-api-webapp"

copy_and_rename_artifact \
    "org.onebusaway:onebusaway-transit-data-federation-webapp:${OBA_VERSION}:war" \
    "onebusaway-transit-data-federation-webapp"

copy_and_rename_artifact \
    "org.onebusaway:onebusaway-transit-data-federation-builder:${OBA_VERSION}:jar:withAllDependencies" \
    "onebusaway-transit-data-federation-builder-withAllDependencies"

copy_and_rename_artifact \
    "com.mysql:mysql-connector-j:${MYSQL_CONNECTOR_VERSION}:jar" \
    "mysql-connector-j"

copy_and_rename_artifact \
    "org.postgresql:postgresql:${POSTGRESQL_CONNECTOR_VERSION}:jar" \
    "postgresql"
