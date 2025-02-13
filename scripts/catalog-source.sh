#!/bin/sh

if [[ -z ${1} ]]; then
    CATALOG_NS="operator-lifecycle-manager"
else
    CATALOG_NS=${1}
fi

CSV=`cat deploy/catalog_resources/redhat/activemq-artemis-operator.v0.8.0.clusterserviceversion.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRD=`cat deploy/crds/broker_v2alpha1_activemqartemis_crd.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRDActivemqartemisaddress=`cat deploy/crds/broker_v2alpha1_activemqartemisaddress_crd.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRDActivemqartemisscaledown=`cat deploy/crds/broker_v2alpha1_activemqartemisscaledown_crd.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
PKG=`cat deploy/catalog_resources/redhat/activemq-artemis.package.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`

cat << EOF > deploy/catalog_resources/redhat/catalog-source.yaml
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ConfigMap
    metadata:
      name: activemq-artemis-resources
      namespace: ${CATALOG_NS}
    data:
      clusterServiceVersions: |
${CSV}
      customResourceDefinitions: |
${CRD}
${CRDActivemqartemisaddress}
${CRDActivemqartemisscaledown}
      packages: >
${PKG}

  - apiVersion: operators.coreos.com/v2alpha1
    kind: CatalogSource
    metadata:
      name: activemq-artemis-resources
      namespace: ${CATALOG_NS}
    spec:
      configMap: activemq-artemis-resources
      displayName: ActiveMQ Artemis Operator
      publisher: Red Hat
      sourceType: internal
    status:
      configMapReference:
        name: activemq-artemis-resources
        namespace: ${CATALOG_NS}
EOF
