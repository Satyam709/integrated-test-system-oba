echo "Building bundle..."
java -Xss4m -Xmx3g \
        -jar ./builder.jar \
        ./gtfs.zip \
        .