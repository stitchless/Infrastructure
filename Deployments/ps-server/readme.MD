This is a local only Repo:

docker run -d -p 5000:5000 --restart=always --name registry registry:2

docker load -i docker-file.tar
docker tag image localhost:5000/image
docker push localhost:5000/image

If new repo is added to kubernetes...
https://wiki.onap.org/display/DW/questions/38118463/how-to-configure-a-kubernetes-cluster-to-use-a-local-insecure-registry

https://medium.com/htc-research-engineering-blog/setup-local-docker-repository-for-local-kubernetes-cluster-354f0730ed3a