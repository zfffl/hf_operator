kubectl-hlf inspect --output ordservice.yaml -o OrdererMSP
kubectl-hlf ca register --name=ord-ca --user=admin --secret=adminpw --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP
kubectl-hlf ca enroll --name=ord-ca --user=admin --secret=adminpw --mspid OrdererMSP --ca-name ca  --output admin-ordservice.yaml
kubectl-hlf utils adduser --userPath=admin-ordservice.yaml --config=ordservice.yaml --username=admin --mspid=OrdererMSP

sleep 2
kubectl-hlf channel generate --output=demo.block --name=demo --organizations Org1MSP --ordererOrganizations OrdererMSP
kubectl-hlf ca enroll --name=ord-ca --namespace=default --user=admin --secret=adminpw --mspid OrdererMSP --ca-name tlsca  --output admin-tls-ordservice.yaml
kubectl-hlf ordnode join --block=demo.block --name=ord-node1 --namespace=default --identity=admin-tls-ordservice.yaml

sleep 2
kubectl-hlf ca register --name=org1-ca --user=admin --secret=adminpw --type=admin  --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP
kubectl-hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP --ca-name ca  --output peer-org1.yaml
kubectl-hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP
kubectl-hlf utils adduser --userPath=peer-org1.yaml --config=org1.yaml --username=admin --mspid=Org1MSP

kubectl-hlf channel join --name=demo --config=org1.yaml --user=admin -p=org1-peer0.default
kubectl-hlf channel inspect --channel=demo --config=org1.yaml --user=admin -p=org1-peer0.default > demo.json
kubectl-hlf channel addanchorpeer --channel=demo --config=org1.yaml --user=admin --peer=org1-peer0.default
rm code.tar.gz asset-transfer-basic-external.tgz
export CHAINCODE_NAME=asset
export CHAINCODE_LABEL=asset
cat << METADATA-EOF > "metadata.json"
{
    "type": "ccaas",
    "label": "${CHAINCODE_LABEL}"
}
METADATA-EOF

cat > "connection.json" <<CONN_EOF
{
  "address": "${CHAINCODE_NAME}:7052",
  "dial_timeout": "10s",
  "tls_required": false
}
CONN_EOF

tar cfz code.tar.gz connection.json
tar cfz asset-transfer-basic-external.tgz metadata.json code.tar.gz
export PACKAGE_ID=$(kubectl-hlf chaincode calculatepackageid --path=asset-transfer-basic-external.tgz --language=node --label=$CHAINCODE_LABEL)
echo "PACKAGE_ID=$PACKAGE_ID"

kubectl-hlf chaincode install --path=./asset-transfer-basic-external.tgz --config=org1.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=org1-peer0.default

sleep 200
kubectl-hlf externalchaincode sync --image=zhangfuli/chaincode-external:latest \
    --name=$CHAINCODE_NAME \
    --namespace=default \
    --package-id=$PACKAGE_ID \
    --tls-required=false \
    --replicas=1
