#!/bin/bash

# Google Cloud プロジェクト設定
PROJECT_ID="passwordgenerator-452012"
SERVICE_NAME="passwordgenerator"
REGION="asia-northeast1"
IMAGE_REGISTRY="asia-northeast1-docker.pkg.dev"
REPOSITORY="cloud-run-source-deploy"

# バージョン管理（デフォルトはタイムスタンプ、引数で指定可能）
VERSION=${1:-"v$(date +%Y%m%d-%H%M%S)"}
IMAGE_URL="${IMAGE_REGISTRY}/${PROJECT_ID}/${REPOSITORY}/${SERVICE_NAME}:${VERSION}"

echo "🚀 Deploying Password Generator to Google Cloud Run..."
echo "📦 Version: $VERSION"
echo "🖼️  Image: $IMAGE_URL"

# プロジェクトを設定
echo "📝 Setting project to $PROJECT_ID"
gcloud config set project $PROJECT_ID

# Cloud Run API の有効化確認
echo "🔧 Enabling required APIs..."
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com

# ステップ1: Docker イメージをビルド
echo "🔨 Building Docker image..."
gcloud builds submit --tag $IMAGE_URL

# ビルド成功確認
if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# ステップ2: Cloud Run にデプロイ
echo "🚢 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_URL \
  --region=$REGION \
  --platform=managed \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1 \
  --max-instances=100 \
  --timeout=300s \
  --concurrency=80 \
  --port=8080 \
  --set-env-vars="GIN_MODE=release"

# デプロイメント成功確認
if [ $? -eq 0 ]; then
    echo "✅ Deployment successful!"
    
    # サービスURLを取得
    SERVICE_URL=$(gcloud run services describe $SERVICE_NAME \
        --region=$REGION \
        --format='value(status.url)')
    
    echo ""
    echo "🎉 Deployment Complete!"
    echo "🌐 Service URL: $SERVICE_URL"
    echo "💡 Health check: $SERVICE_URL/health"
    echo "📦 Image: $IMAGE_URL"
    echo ""
    
    # 簡単なヘルスチェック
    echo "🔍 Performing health check..."
    sleep 5  # サービス起動を待つ
    
    if curl -f -s "$SERVICE_URL/health" > /dev/null; then
        echo "✅ Health check passed!"
        echo "🚀 Service is ready!"
    else
        echo "⚠️  Health check failed. Service might still be starting up..."
        echo "🔧 Try accessing: $SERVICE_URL"
    fi
    
    echo ""
    echo "📝 Usage:"
    echo "  ./deploy.sh              # Deploy with auto-generated version"
    echo "  ./deploy.sh v1.0.0       # Deploy with specific version"
    echo "  ./deploy.sh latest       # Deploy with 'latest' tag"
    
else
    echo "❌ Deployment failed!"
    exit 1
fi