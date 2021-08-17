helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

helm install dev-mysql --namespace dev \
  --set auth.rootPassword=asdfasdfkjhuh234iuy \
  --set auth.database=base_database \
  --set auth.username=stitch \
  --set auth.password=sdfh78sdf8sdf \
  --set primary.service.type=LoadBalancer \
  --set secondary.persistence.selector=mysql-pv-volume \
  bitnami/mysql

