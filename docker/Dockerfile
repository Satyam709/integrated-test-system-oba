FROM alpine:latest

# Set environment variables
ENV CATALINA_HOME="/usr/local/tomcat"
ENV CATALINA_OPTS="-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005"
ENV TZ="America/Los_Angeles"


# Install essential runtime dependencies first
RUN apk add --no-cache \
    openjdk11 \
    curl \
    git \
    jq \
    python3 \
    py3-pip \
    supervisor \
    tzdata \
    unzip \
    vim \
    xmlstarlet \
    bash

# Setup java
ENV JAVA_HOME="/usr/lib/jvm/java-11-openjdk"
ENV PATH="$JAVA_HOME/bin:$PATH"

RUN mkdir -p /root/.m2/repository

# Download and install Tomcat
RUN mkdir -p ${CATALINA_HOME} && \
    curl -o /tmp/apache-tomcat.tar.gz https://archive.apache.org/dist/tomcat/tomcat-8/v8.5.100/bin/apache-tomcat-8.5.100.tar.gz && \
    tar -xzf /tmp/apache-tomcat.tar.gz -C ${CATALINA_HOME} --strip-components=1 && \
    rm -f /tmp/apache-tomcat.tar.gz

# Move the gtfs static_tests file
COPY ./testdata/gtfs.zip \
    /bundle/gtfs.zip

# Build bundle
WORKDIR /bundle
#COPY ./docker/build_bundle.sh .
#COPY ./docker/copy_resources.sh .
COPY ./oba-artifacts/onebusaway-transit-data-federation-builder-withAllDependencies.jar ./builder.jar
COPY ./testdata/gtfs.zip .
COPY ./docker/scripts/build_bundle.sh .
RUN chmod +x build_bundle.sh
RUN ./build_bundle.sh

# Set the configured time zone
RUN ln -fs /usr/share/zoneinfo/$TZ /etc/localtime

# Set up the host-manager and manager webapps
COPY ./docker/config/tomcat-users.xml $CATALINA_HOME/conf/
COPY ./docker/config/host-manager_context.xml $CATALINA_HOME/webapps/docs/META-INF/context.xml
COPY ./docker/config/host-manager_context.xml $CATALINA_HOME/webapps/examples/META-INF/context.xml
COPY ./docker/config/host-manager_context.xml $CATALINA_HOME/webapps/host-manager/META-INF/context.xml
COPY ./docker/config/host-manager_context.xml $CATALINA_HOME/webapps/manager/META-INF/context.xml
COPY ./docker/config/host-manager_context.xml $CATALINA_HOME/webapps/ROOT/META-INF/context.xml


# Copy the OBA WAR file and configuration

# transit-data-federation-webapp
COPY ./oba-artifacts/onebusaway-transit-data-federation-webapp.war $CATALINA_HOME/webapps
RUN unzip -o $CATALINA_HOME/webapps/onebusaway-transit-data-federation-webapp.war -d $CATALINA_HOME/webapps/onebusaway-transit-data-federation-webapp
RUN rm -rf $CATALINA_HOME/webapps/onebusaway-transit-data-federation-webapp.war
COPY ./docker/config/onebusaway-transit-data-federation-webapp-data-sources.xml $CATALINA_HOME/webapps/onebusaway-transit-data-federation-webapp/WEB-INF/classes/data-sources.xml

# onebusaway-api-webapp
COPY ./oba-artifacts/onebusaway-api-webapp.war $CATALINA_HOME/webapps
RUN unzip -o $CATALINA_HOME/webapps/onebusaway-api-webapp.war -d $CATALINA_HOME/webapps/onebusaway-api-webapp
RUN rm -rf $CATALINA_HOME/webapps/onebusaway-api-webapp.war
COPY ./docker/config/onebusaway-api-webapp-data-sources.xml $CATALINA_HOME/webapps/onebusaway-api-webapp/WEB-INF/classes/data-sources.xml


# Setup libfaketime
WORKDIR /libfaketimeInstall
RUN git clone https://github.com/wolfcw/libfaketime.git

RUN apk add --no-cache --virtual .build-deps \
    gcc \
    musl-dev \
    make

# Build and install libfaketime
WORKDIR /libfaketimeInstall/libfaketime
RUN make && make install

# Remove build dependencies after installing libfaketime to reduce image size
RUN apk del .build-deps

# install gcompat
RUN apk add gcompat
RUN apk add --no-cache libc6-compat
# Set libfaketime environment variables
ENV LD_PRELOAD=/usr/local/lib/faketime/libfaketime.so.1
ENV FAKETIME_DONT_FAKE_MONOTONIC=1
# create faketimerc file
RUN touch /etc/faketimerc

# MySQL Connector
WORKDIR $CATALINA_HOME/lib
COPY ./oba-artifacts/mysql-connector-j.jar .

WORKDIR /

# Expose an additional port for debugging
EXPOSE 5005

# Start Tomcat
# CMD ["sh", "-c", "echo \"Starting Tomcat...\" && $CATALINA_HOME/bin/catalina.sh run"]