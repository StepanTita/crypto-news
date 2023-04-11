# Build a docker container and
# push it to Docker Hub so that it can
# be deployed to EC2.

docker login -u "$IMAGE_REPO" -p "$DOCKERHUB_TOKEN"

services=("migrator" "parser" "telegram-bot" "twitter-bot")
for IMAGE_NAME in "${services[@]}"; do
  echo "Processing $IMAGE_NAME..."

  docker build -t "$IMAGE_REPO"/"$IMAGE_NAME":"$IMAGE_TAG" . -f ./"$IMAGE_NAME"/Dockerfile
  docker push "$IMAGE_REPO"/"$IMAGE_NAME":"$IMAGE_TAG"
  docker tag "$IMAGE_REPO"/"$IMAGE_NAME":"$IMAGE_TAG" "$IMAGE_REPO"/"$IMAGE_NAME":latest
  docker push "$IMAGE_REPO"/"$IMAGE_NAME":latest

  echo "Pushed" "$IMAGE_REPO"/"$IMAGE_NAME":"$IMAGE_TAG"
  echo "Pushed" "$IMAGE_REPO"/"$IMAGE_NAME":latest
done