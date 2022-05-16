# 循环10次创建单orderer、单org（1个peer）

for i in {1..1}; do
  rm -rf ./out.txt

  kubectl delete fabricorderernodes.hlf.kungfusoftware.es --all-namespaces --all

  if [ $? -ne 0 ]; then
    echo "delete orderer failed"
  else
    echo "delete orderer success"
    kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
    if [ $? -ne 0 ]; then
      echo "delete peers failed"
    else
      echo "delete peers success"
      kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
      if [ $? -ne 0 ]; then
        echo "delete ca failed"
      else
        echo "delete ca success"
      fi
    fi
  fi

  sleep 60
  echo $i

  time1=$(date)
  echo $time1
  kubectl-hlf ca create --storage-class=standard --capacity=0.2Gi --name=org1-ca --enroll-id=enroll --enroll-pw=enrollpw
  kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
  if [ $? -ne 0 ]; then
    echo "peer ca failed"
  else
    echo "peer ca succeed"
    time2=$(date)
    echo $time2
    kubectl-hlf ca register --name=org1-ca --user=peer --secret=peerpw --type=peer --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP
    kubectl-hlf peer create --statedb=leveldb --storage-class=standard --enroll-id=peer --mspid=Org1MSP --enroll-pw=peerpw --capacity=0.5Gi --name=org1-peer0 --ca-name=org1-ca.default
    kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all
    if [ $? -ne 0 ]; then
      echo "peer failed"
    else
      echo "peer succeed"
      time3=$(date)
      echo $time3
      kubectl-hlf ca create --storage-class=standard --capacity=0.2Gi --name=ord-ca --enroll-id=enroll --enroll-pw=enrollpw
      kubectl-wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
      if [ $? -ne 0 ]; then
        echo "orderer ca failed"
      else
        echo "orderer ca success"
        time4=$(date)
        echo $time4
        kubectl-hlf ca register --name=ord-ca --user=orderer --secret=ordererpw --type=orderer --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP
        kubectl-hlf ordnode create --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP --enroll-pw=ordererpw --capacity=0.2Gi --name=ord-node1 --ca-name=ord-ca.default
        kubectl wait --timeout=180s --for=condition=Running fabricorderernodes.hlf.kungfusoftware.es --all
        if [ $? -ne 0 ]; then
          echo "orderer failed"
        else
          echo "orderer success"
          time5=$(date)
          echo $time5
        fi
      fi
    fi
  fi

done
