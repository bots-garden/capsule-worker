#!/bin/bash

TAG="v0.0.0"

git add .
git commit -m "📦 updates modules for ${TAG}"

git tag ${TAG}
git tag commons/${TAG}

git push origin main ${TAG} commons/${TAG}

