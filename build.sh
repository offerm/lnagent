IMAGE=lnagent
echo "Building container... $IMAGE"
docker build . -f Dockerfile --rm=true  -t $IMAGE
